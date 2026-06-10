import { useForm, useSelect } from '@refinedev/antd';
import { Create } from '@refinedev/antd';
import { Form, Input, InputNumber, Select, Switch, Space, Divider } from 'antd';
import { rublesToKopeks } from '../../utils/api';
import { ImageUpload } from '../../components/ImageUpload';

export const ProductCreate = () => {
  const { formProps, saveButtonProps } = useForm();

  const { selectProps: categorySelectProps } = useSelect({
    resource: 'categories',
    optionLabel: 'name',
    optionValue: 'id',
  });

  const { selectProps: firmSelectProps } = useSelect({
    resource: 'firms',
    optionLabel: 'name',
    optionValue: 'id',
  });

  const { selectProps: tagSelectProps } = useSelect({
    resource: 'tags',
    optionLabel: 'name',
    optionValue: 'id',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form
        {...formProps}
        layout="vertical"
        onFinish={(values) => {
          const priceKopeks = rublesToKopeks(values.price_rubles);
          const { price_rubles, images, ...rest } = values;
          formProps.onFinish?.({
            ...rest,
            price_kopeks: priceKopeks,
            image_urls: images || [],
          });
        }}
      >
        <Divider>Основная информация</Divider>
        <Form.Item
          label="Название товара"
          name="name"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Например: iPhone 15 Pro" />
        </Form.Item>

        <Form.Item
          label="Описание"
          name="description"
        >
          <Input.TextArea rows={4} placeholder="Подробное описание товара" />
        </Form.Item>

        <Space style={{ width: '100%' }} size="large">
          <Form.Item
            label="Категория"
            name="category_id"
            rules={[{ required: true, message: 'Выберите категорию' }]}
            style={{ width: 300 }}
          >
            <Select {...categorySelectProps} placeholder="Выберите категорию" />
          </Form.Item>

          <Form.Item
            label="Бренд"
            name="firm_id"
            style={{ width: 300 }}
          >
            <Select {...firmSelectProps} placeholder="Выберите бренд" allowClear />
          </Form.Item>
        </Space>

        <Space style={{ width: '100%' }} size="large">
          <Form.Item
            label="Цена (Br)"
            name="price_rubles"
            rules={[{ required: true, message: 'Введите цену' }]}
          >
            <InputNumber
              min={0}
              step={0.01}
              precision={2}
              placeholder="0.00"
              style={{ width: 200 }}
            />
          </Form.Item>

          <Form.Item
            label="Количество на складе"
            name="stock_quantity"
            rules={[{ required: true, message: 'Введите количество' }]}
          >
            <InputNumber
              min={0}
              placeholder="0"
              style={{ width: 150 }}
            />
          </Form.Item>
        </Space>

        <Form.Item
          label="Теги"
          name="tag_ids"
        >
          <Select
            {...tagSelectProps}
            mode="multiple"
            placeholder="Выберите теги"
            style={{ width: '100%' }}
          />
        </Form.Item>

        <Form.Item
          label="Активен"
          name="is_active"
          valuePropName="checked"
          initialValue={true}
        >
          <Switch />
        </Form.Item>

        <Divider>Изображения</Divider>
        
        <Form.Item
          label="Изображения товара"
          name="images"
          help="Можно загрузить до 5 изображений. Первое изображение будет главным."
        >
          <ImageUpload maxCount={5} />
        </Form.Item>
      </Form>
    </Create>
  );
};
