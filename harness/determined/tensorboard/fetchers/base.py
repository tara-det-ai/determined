import abc
from typing import Any, Callable, Dict, List


class Fetcher(metaclass=abc.ABCMeta):
    """Abstract base class for Tensorflow fetchers.

    Responsible for syncing tensorflow files from a list of storage paths to a local directory.
    """

    storage_paths: List[str]

    @abc.abstractmethod
    def __init__(self, storage_config: Dict[str, Any], storage_paths: List[str], local_dir: str):
        pass

    @abc.abstractmethod
    def fetch_new(self, new_file_callback: Callable = lambda: None) -> int:
        """Fetches changed files found in storage paths to local disk.

        Arguments:
            new_file_callback (Callable, optional): Callback function that
                is fired each time a new file is fetched

        Returns: count of new files fetched.

        """
        pass
