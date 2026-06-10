import { useLogin } from '@refinedev/core';
import { Form, Input, Button, Card, Typography, Alert } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useState } from 'react';

const { Title } = Typography;

export const Login = () => {
  const { mutate: login } = useLogin();
  const [error, setError] = useState<string | null>(null);

  const onFinish = (values: { username: string; password: string }) => {
    setError(null);
    login(values, {
      onError: (error) => {
        setError(error.message || 'Login failed');
      },
    });
  };

  return (
    <div style={{
      height: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <Card style={{ width: 400, boxShadow: '0 4px 12px rgba(0,0,0,0.15)' }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={2} style={{ marginBottom: 8 }}>AlphaSales Admin</Title>
          <Typography.Text type="secondary">Войдите в панель управления</Typography.Text>
        </div>

        {error && (
          <Alert
            message={error}
            type="error"
            showIcon
            closable
            onClose={() => setError(null)}
            style={{ marginBottom: 16 }}
          />
        )}

        <Form
          name="login"
          onFinish={onFinish}
          layout="vertical"
          requiredMark={false}
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: 'Введите имя пользователя' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="Имя пользователя"
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: 'Введите пароль' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Пароль"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" size="large" block>
              Войти
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};
