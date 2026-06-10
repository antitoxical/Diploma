import { useForm, useSelect } from '@refinedev/antd';
import { Edit } from '@refinedev/antd';
import { Form, Input, Select, Switch } from 'antd';
import { ImageUpload } from '../../components/ImageUpload';

export const CategoryEdit = () => {
  const { formProps, saveButtonProps, queryResult } = useForm({
    redirect: 'list',
  });

  const categoryData = queryResult?.data?.data;

  const { selectProps: parentCategorySelectProps } = useSelect({
    resource: 'categories',
    optionLabel: 'name',
    optionValue: 'id',
    defaultValue: categoryData?.parent_id,
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form
        {...formProps}
        layout="vertical"
        onFinish={(values) => {
          const { image, image_url: urlFromField, ...rest } = values;
          const uploaded = image && image.length > 0 ? image[0] : undefined;
          const image_url =
            uploaded ||
            (typeof urlFromField === 'string' && urlFromField.trim() ? urlFromField.trim() : undefined) ||
            categoryData?.image_url ||
            null;
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
          help="Оставьте пустым для корневой категории"
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
          getValueFromEvent={(value) => value}
          getValueProps={(value) => {
            if (categoryData?.image_url) {
              return {
                value: [categoryData.image_url],
              };
            }
            return { value };
          }}
        >
          <ImageUpload maxCount={1} />
        </Form.Item>

        <Form.Item
          label="Активна"
          name="is_active"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
