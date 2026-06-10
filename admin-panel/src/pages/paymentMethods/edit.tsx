import { Edit, useForm } from '@refinedev/antd';
import { Form, Input, Switch } from 'antd';

export const PaymentMethodEdit = () => {
  const { formProps, saveButtonProps } = useForm({ redirect: 'list' });

  return (
    <Edit saveButtonProps={saveButtonProps} title="Изменить способ оплаты">
      <Form {...formProps} layout="vertical">
        <Form.Item label="Название" name="name" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item label="Описание" name="description">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item label="Активен" name="is_active" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
