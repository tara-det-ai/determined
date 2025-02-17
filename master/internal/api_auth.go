package internal

import (
	"context"
	"crypto/sha512"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/determined-ai/determined/master/internal/command"
	"github.com/determined-ai/determined/master/internal/db"
	"github.com/determined-ai/determined/master/internal/grpcutil"
	"github.com/determined-ai/determined/master/internal/user"
	"github.com/determined-ai/determined/master/pkg/model"
	"github.com/determined-ai/determined/proto/pkg/apiv1"
)

const clientSidePasswordSalt = "GubPEmmotfiK9TMD6Zdw" // #nosec G101

// replicateClientSideSaltAndHash replicates the password salt and hash done on the client side.
// We need this because we hash passwords on the client side, but when SCIM posts a user with
// a password to password sync, it doesn't - so when we try to log in later, we get a weird,
// unrecognizable sha512 hash from the frontend.
func replicateClientSideSaltAndHash(password string) string {
	if password == "" {
		return password
	}
	sum := sha512.Sum512([]byte(clientSidePasswordSalt + password))
	return fmt.Sprintf("%x", sum)
}

func (a *apiServer) Login(
	ctx context.Context, req *apiv1.LoginRequest,
) (*apiv1.LoginResponse, error) {
	if a.m.config.InternalConfig.ExternalSessions.JwtKey != "" {
		return nil, status.Error(codes.FailedPrecondition, "authentication is configured to be external")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument: username")
	}

	userModel, err := user.UserByUsername(req.Username)
	switch err {
	case nil:
	case db.ErrNotFound:
		return nil, grpcutil.ErrInvalidCredentials
	default:
		return nil, err
	}

	var hashedPassword string
	if req.IsHashed {
		hashedPassword = req.Password
	} else {
		hashedPassword = replicateClientSideSaltAndHash(req.Password)
	}

	if !userModel.ValidatePassword(hashedPassword) {
		return nil, grpcutil.ErrInvalidCredentials
	}

	if !userModel.Active {
		return nil, grpcutil.ErrNotActive
	}
	token, err := a.m.db.StartUserSession(userModel)
	if err != nil {
		return nil, err
	}
	fullUser, err := getUser(a.m.db, userModel.ID)
	return &apiv1.LoginResponse{Token: token, User: fullUser}, err
}

func (a *apiServer) CurrentUser(
	ctx context.Context, _ *apiv1.CurrentUserRequest,
) (*apiv1.CurrentUserResponse, error) {
	user, _, err := grpcutil.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	fullUser, err := getUser(a.m.db, user.ID)
	return &apiv1.CurrentUserResponse{User: fullUser}, err
}

func (a *apiServer) Logout(
	ctx context.Context, _ *apiv1.LogoutRequest,
) (*apiv1.LogoutResponse, error) {
	_, userSession, err := grpcutil.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if userSession == nil {
		return nil, status.Error(codes.InvalidArgument,
			"cannot manually logout of an allocation session")
	}

	err = a.m.db.DeleteUserSessionByID(userSession.ID)
	return &apiv1.LogoutResponse{}, err
}

func redirectToLogin(c echo.Context) error {
	return c.Redirect(
		http.StatusSeeOther,
		fmt.Sprintf("/det/login?redirect=%s", c.Request().URL),
	)
}

// processProxyAuthentication is a middleware processing function that attempts
// to authenticate incoming HTTP requests coming through proxies.
func processProxyAuthentication(c echo.Context) (done bool, err error) {
	user, _, err := user.GetService().UserAndSessionFromRequest(c.Request())
	if errors.Is(err, db.ErrNotFound) {
		return true, redirectToLogin(c)
	} else if err != nil {
		return true, err
	}
	if !user.Active {
		return true, redirectToLogin(c)
	}

	taskID := model.TaskID(c.Param("service"))
	var ctx context.Context

	if c.Request() == nil || c.Request().Context() == nil {
		ctx = context.TODO()
	} else {
		ctx = c.Request().Context()
	}

	spec, err := db.IdentifyTask(ctx, taskID)
	if err != nil {
		return true, err
	}

	var ok bool
	if spec.TaskType == model.TaskTypeTensorboard {
		ok, err = command.AuthZProvider.Get().CanGetTensorboard(
			ctx, *user, spec.WorkspaceID, spec.ExperimentIDs, spec.TrialIDs)
	} else {
		ok, err = command.AuthZProvider.Get().CanGetNSC(
			ctx, *user, spec.WorkspaceID)
	}
	if err != nil {
		return true, err
	} else if !ok {
		return true, echo.NewHTTPError(http.StatusNotFound, "service not found: "+taskID)
	}

	return false, nil
}
