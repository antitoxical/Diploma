import { useForm } from '@refinedev/antd';
import { Create } from '@refinedev/antd';
import { Form, Input, Switch } from 'antd';

export const AdminCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Имя пользователя"
          name="username"
          rules={[
            { required: true, message: 'Введите имя пользователя' },
            { min: 3, message: 'Минимум 3 символа' },
          ]}
        >
          <Input placeholder="admin" />
        </Form.Item>

        <Form.Item
          label="Пароль"
          name="password"
          rules={[
            { required: true, message: 'Введите пароль' },
            { min: 6, message: 'Минимум 6 символов' },
          ]}
        >
          <Input.Password placeholder="Введите пароль" />
        </Form.Item>

        <Form.Item
          label="Подтверждение пароля"
          name="password_confirm"
          dependencies={['password']}
          rules={[
            { required: true, message: 'Подтвердите пароль' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('Пароли не совпадают!'));
              },
            }),
          ]}
        >
          <Input.Password placeholder="Подтвердите пароль" />
        </Form.Item>

        <Form.Item
          label="Активен"
          name="is_active"
          valuePropName="checked"
          initialValue={true}
        >
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
