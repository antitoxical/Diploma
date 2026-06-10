import { useForm } from '@refinedev/antd';
import { Edit } from '@refinedev/antd';
import { Form, Input } from 'antd';

export const FirmEdit = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Название бренда"
          name="name"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Например: Apple" />
        </Form.Item>
      </Form>
    </Edit>
  );
};
