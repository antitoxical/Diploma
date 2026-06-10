import { useForm } from '@refinedev/antd';
import { Create } from '@refinedev/antd';
import { Form, Input } from 'antd';

export const FirmCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Название бренда"
          name="name"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Например: Apple" />
        </Form.Item>
      </Form>
    </Create>
  );
};
