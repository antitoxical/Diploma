import { useForm } from '@refinedev/antd';
import { Create } from '@refinedev/antd';
import { Form, Input, InputNumber, Switch, DatePicker } from 'antd';

export const PromocodeCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    redirect: 'list',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Код промокода"
          name="code"
          rules={[
            { required: true, message: 'Введите код' },
            { pattern: /^[A-Z0-9]+$/, message: 'Используйте только заглавные буквы и цифры' }
          ]}
        >
          <Input placeholder="SUMMER2026" style={{ textTransform: 'uppercase' }} />
        </Form.Item>

        <Form.Item
          label="Размер скидки (%)"
          name="discount_percent"
          rules={[{ required: true, message: 'Введите размер скидки' }]}
        >
          <InputNumber
            min={1}
            max={100}
            placeholder="10"
            style={{ width: 200 }}
          />
        </Form.Item>

        <Form.Item
          label="Максимальное количество использований"
          name="max_uses"
          help="Оставьте пустым для неограниченного использования"
        >
          <InputNumber
            min={1}
            placeholder="100"
            style={{ width: 200 }}
          />
        </Form.Item>

        <Form.Item
          label="Дата истечения"
          name="expires_at"
          help="Оставьте пустым для промокода без срока действия"
        >
          <DatePicker
            showTime
            format="DD.MM.YYYY HH:mm"
            style={{ width: 300 }}
            placeholder="Выберите дату"
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
      </Form>
    </Create>
  );
};
