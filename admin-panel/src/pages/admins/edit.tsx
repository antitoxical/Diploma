import { useForm } from '@refinedev/antd';
import { Edit } from '@refinedev/antd';
import { Form, Input, Switch, Alert } from 'antd';

export const AdminEdit = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Alert
          message="Внимание"
          description="Оставьте поле пароля пустым, если не хотите его изменять."
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
        />

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
          label="Новый пароль"
          name="password"
          help="Оставьте пустым, чтобы не менять пароль"
          rules={[
            { min: 6, message: 'Минимум 6 символов' },
          ]}
        >
          <Input.Password placeholder="Введите новый пароль (опционально)" />
        </Form.Item>

        <Form.Item
          label="Подтверждение нового пароля"
          name="password_confirm"
          dependencies={['password']}
          rules={[
            ({ getFieldValue }) => ({
              validator(_, value) {
                const password = getFieldValue('password');
                if (!password || !value || password === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('Пароли не совпадают!'));
              },
            }),
          ]}
        >
          <Input.Password placeholder="Подтвердите новый пароль" />
        </Form.Item>

        <Form.Item
          label="Активен"
          name="is_active"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
