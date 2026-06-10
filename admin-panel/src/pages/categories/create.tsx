import { useForm, useSelect } from '@refinedev/antd';
import { Create } from '@refinedev/antd';
import { Form, Input, Select, Switch } from 'antd';
import { ImageUpload } from '../../components/ImageUpload';

export const CategoryCreate = () => {
  const { formProps, saveButtonProps } = useForm();

  const { selectProps: parentCategorySelectProps } = useSelect({
    resource: 'categories',
    optionLabel: 'name',
    optionValue: 'id',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form
        {...formProps}
        layout="vertical"
        onFinish={(values) => {
          const { image, image_url: urlFromField, ...rest } = values;
          const uploaded = image && image.length > 0 ? image[0] : undefined;
          const image_url =
            uploaded ||
            (typeof urlFromField === 'string' && urlFromField.trim() ? urlFromField.trim() : null);
          formProps.onFinish?.({
            ...rest,
            image_url,
          });
        }}
      >
        <Form.Item
          label="Название категории"
          name="name"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Например: Смартфоны" />
        </Form.Item>

        <Form.Item
          label="Родительская категория"
          name="parent_id"
          help="Оставьте пустым для создания корневой категории"
        >
          <Select
            {...parentCategorySelectProps}
            placeholder="Выберите родительскую категорию"
            allowClear
          />
        </Form.Item>

        <Form.Item
          label="URL изображения"
          name="image_url"
        >
          <Input placeholder="https://example.com/image.jpg" />
        </Form.Item>

        <Form.Item
          label="Изображение категории"
          name="image"
        >
          <ImageUpload maxCount={1} />
        </Form.Item>

        <Form.Item
          label="Активна"
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
