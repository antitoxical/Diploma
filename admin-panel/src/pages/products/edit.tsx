import { useForm, useSelect } from '@refinedev/antd';

import { Edit } from '@refinedev/antd';

import { Form, Input, InputNumber, Select, Switch, Space, Divider } from 'antd';

import { ImageUpload } from '../../components/ImageUpload';

import type { Product } from '../../types';

import {

  mapProductToFormValues,

  productFormValuesToUpdatePayload,

  type ProductFormValues,

} from '../../utils/productForm';



export const ProductEdit = () => {

  const { formProps, saveButtonProps, queryResult } = useForm<Product, unknown, ProductFormValues>({

    redirect: 'list',

    queryOptions: {

      select: (response) => ({

        ...response,

        data: mapProductToFormValues(response.data as Product),

      }),

    },

  });



  const formRecord = queryResult?.data?.data as ProductFormValues | undefined;



  const { selectProps: categorySelectProps } = useSelect({

    resource: 'categories',

    optionLabel: 'name',

    optionValue: 'id',

    defaultValue: formRecord?.category_id,

  });



  const { selectProps: firmSelectProps } = useSelect({

    resource: 'firms',

    optionLabel: 'name',

    optionValue: 'id',

    defaultValue: formRecord?.firm_id,

  });



  const { selectProps: tagSelectProps } = useSelect({

    resource: 'tags',

    optionLabel: 'name',

    optionValue: 'id',

  });



  return (

    <Edit saveButtonProps={saveButtonProps}>

      <Form

        {...formProps}

        layout="vertical"

        onFinish={(values: ProductFormValues) => {

          formProps.onFinish?.(productFormValuesToUpdatePayload(values, formRecord));

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

    </Edit>

  );

};


