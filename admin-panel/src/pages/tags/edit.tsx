import { useForm } from '@refinedev/antd';
import { Edit } from '@refinedev/antd';
import { Form, Input } from 'antd';

export const TagEdit = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Название тега"
          name="name"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Например: Новинка" />
        </Form.Item>

        <Form.Item
          label="Цвет (HEX)"
          name="color"
          rules={[
            { required: true, message: 'Введите цвет' },
            { pattern: /^#[0-9A-F]{6}$/i, message: 'Неверный формат цвета (используйте #RRGGBB)' }
          ]}
        >
          <Input
            placeholder="#22c55e"
            addonBefore={
              <div style={{ width: 20, height: 20, backgroundColor: formProps?.form?.getFieldValue('color') || '#22c55e', borderRadius: 4 }} />
            }
          />
        </Form.Item>
      </Form>
    </Edit>
  );
};
