import {
  Form as AntdForm,
  FormInstance as AntdFormInstance,
  FormItemProps as AntdFormItemProps,
} from 'antd';
import { FormListFieldData as AntdFormListFieldData } from 'antd/lib/form/FormList';
import { FieldData as AntdFieldData, NamePath as AntdNamePath } from 'rc-field-form/lib/interface';
import React, { FC, ReactNode, Ref } from 'react';

import { Primitive } from 'shared/types';

type Rules = AntdFormItemProps['rules']; // https://github.com/ant-design/ant-design/issues/39466
type GridCol = {
  span: number;
};
type TriggerEvent = 'onChange' | 'onSubmit';

interface FormItemProps {
  children?: ReactNode;
  className?: string;
  dependencies?: AntdNamePath[];
  extra?: ReactNode;
  field?: AntdFormListFieldData;
  hidden?: boolean;
  initialValue?: string | number | Primitive;
  label?: ReactNode;
  labelCol?: GridCol; // https://ant.design/components/grid#col
  max?: number;
  maxMessage?: string;
  name?: string | number | (string | number)[];
  noStyle?: boolean;
  required?: boolean;
  requiredMessage?: string;
  rules?: Rules; // https://ant.design/components/form#rule
  validateMessage?: string;
  validateStatus?: 'success' | 'warning' | 'error' | 'validating';
  validateTrigger?: TriggerEvent[];
  valuePropName?: string;
}

const FormItem: React.FC<FormItemProps> = ({
  children,
  label,
  labelCol = { span: 24 },
  max,
  maxMessage,
  required,
  requiredMessage,
  rules = [],
  validateMessage,
  ...props
}: FormItemProps) => {
  if (required) rules.push({ message: requiredMessage || `${label} required`, required: true });
  if (max) rules.push({ max, message: maxMessage || `${label} cannot exceed ${max} characters` });

  return (
    <AntdForm.Item
      help={validateMessage}
      label={label}
      labelCol={labelCol}
      required={required}
      rules={rules}
      {...props}>
      {children}
    </AntdForm.Item>
  );
};

interface FormProps {
  autoComplete?: string;
  children?: ReactNode;
  className?: string;
  fields?: AntdFieldData[];
  form?: AntdFormInstance;
  hidden?: boolean;
  id?: string;
  initialValues?: object;
  labelCol?: GridCol;
  layout?: 'horizontal' | 'vertical' | 'inline';
  name?: string;
  onFieldsChange?: (changedFields: any, allFields: any) => void;
  onFinish?: (values: any) => void;
  onValuesChange?: (changedValues: any, allValues: any) => void;
  ref?: Ref<AntdFormInstance>;
  wrapperCol?: GridCol;
}

type Form = JSX.Element & {
  Item?: FC<FormItemProps>;
  List?: typeof AntdForm.List;
  useForm?: typeof AntdForm.useForm;
};

const Form = (props: FormProps): JSX.Element => {
  return <AntdForm {...props} />;
};

Form.Item = FormItem;
Form.List = AntdForm.List;
Form.useForm = AntdForm.useForm;
Form.useWatch = AntdForm.useWatch;

export type FormInstance<Values = any> = AntdFormInstance<Values>;

export type FormListFieldData = AntdFormListFieldData;

export default Form;
