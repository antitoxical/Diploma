import { Create, useForm } from '@refinedev/antd';
import { Form, Input, Switch, InputNumber, Select } from 'antd';
import { DELIVERY_KIND_LABELS, type DeliveryMethodKind } from '../../types';

const kindOptions = (Object.keys(DELIVERY_KIND_LABELS) as DeliveryMethodKind[]).map(k => ({
  value: k,
  label: DELIVERY_KIND_LABELS[k],
}));

export const DeliveryMethodCreate = () => {
  const { formProps, saveButtonProps } = useForm({ redirect: 'list' });

  return (
    <Create saveButtonProps={saveButtonProps} title="Добавить способ доставки">
      <Form {...formProps} layout="vertical" initialValues={{ kind: 'courier', is_active: true }}>
        <Form.Item label="Название" name="name" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item
          label="Тип (поля на оформлении заказа)"
          name="kind"
          rules={[{ required: true }]}
          extra="Курьер — улица и дом; Самовывоз — только город; СДЭК — город и адрес ПВЗ."
        >
          <Select options={kindOptions} />
        </Form.Item>
        <Form.Item label="Цена (копейки)" name="price_kopeks" rules={[{ required: true }]}>
          <InputNumber min={0} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item label="Описание" name="description">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item label="Активен" name="is_active" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
