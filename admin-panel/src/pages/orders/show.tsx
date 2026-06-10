import { useOne } from '@refinedev/core';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card, Descriptions, Tag, Table, Button, Select, Space,
  Typography, Divider, message, Spin,
} from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { axiosInstance, formatPrice } from '../../utils/api';
import { useState } from 'react';
import dayjs from 'dayjs';

const { Title } = Typography;

const STATUS_COLORS: Record<string, string> = {
  created: 'blue',
  paid: 'green',
  delivering: 'orange',
  completed: 'success',
  cancelled: 'red',
};

const STATUS_LABELS: Record<string, string> = {
  created: 'Создан',
  paid: 'Оплачен',
  delivering: 'Доставляется',
  completed: 'Завершён',
  cancelled: 'Отменён',
};

const STATUS_OPTIONS = Object.entries(STATUS_LABELS).map(([value, label]) => ({ value, label }));

export const OrderShow = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [newStatus, setNewStatus] = useState<string>('');
  const [comment, setComment] = useState('');
  const [updating, setUpdating] = useState(false);

  const { data, isLoading, refetch } = useOne({
    resource: 'orders',
    id: id!,
  });

  const order = data?.data as any;

  const handleStatusUpdate = async () => {
    if (!newStatus) { message.warning('Выберите новый статус'); return; }
    setUpdating(true);
    try {
      await axiosInstance.patch(`/admin/orders/${id}/status`, { status: newStatus, comment });
      message.success('Статус обновлён');
      setNewStatus('');
      setComment('');
      refetch();
    } catch (e: any) {
      message.error(e?.response?.data?.data || 'Ошибка при обновлении статуса');
    } finally {
      setUpdating(false);
    }
  };

  if (isLoading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />;
  if (!order) return <div>Заказ не найден</div>;

  const products: any[] = order.products ?? [];

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/orders')}>
          Назад к заказам
        </Button>
      </Space>

      <Title level={3}>Заказ #{String(order.id).substring(0, 8)}…</Title>

      <Card style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID заказа" span={2}>
            <code>{order.id}</code>
          </Descriptions.Item>
          <Descriptions.Item label="Пользователь">
            <code>{order.user_id}</code>
          </Descriptions.Item>
          <Descriptions.Item label="Telegram ID">
            {order.telegram_id ?? '—'}
          </Descriptions.Item>
          <Descriptions.Item label="Статус">
            <Tag color={STATUS_COLORS[order.status]}>{STATUS_LABELS[order.status] ?? order.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Дата создания">
            {dayjs(order.created_at).format('DD.MM.YYYY HH:mm')}
          </Descriptions.Item>
          <Descriptions.Item label="Сумма">
            {formatPrice(order.total_amount_kopeks)}
          </Descriptions.Item>
          <Descriptions.Item label="Адрес доставки">
            {order.address_snapshot ?? order.delivery_address ?? '—'}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Товары в заказе" style={{ marginBottom: 16 }}>
        <Table
          dataSource={products}
          rowKey="id"
          pagination={false}
          columns={[
            {
              title: 'Товар',
              dataIndex: 'product_name',
              render: (name: string, r: any) =>
                name?.trim() ? name : <code>{r.product_id?.substring(0, 8)}…</code>,
            },
            { title: 'Кол-во', dataIndex: 'quantity' },
            { title: 'Цена за ед.', dataIndex: 'price_at_purchase_kopeks', render: (v: number) => formatPrice(v) },
            { title: 'Итого', render: (_: any, r: any) => formatPrice(r.quantity * r.price_at_purchase_kopeks) },
          ]}
        />
      </Card>

      <Card title="Изменить статус заказа">
        <Space direction="vertical" style={{ width: '100%' }}>
          <Select
            style={{ width: 220 }}
            placeholder="Выберите новый статус"
            value={newStatus || undefined}
            options={STATUS_OPTIONS}
            onChange={setNewStatus}
          />
          <input
            style={{ width: '100%', padding: '4px 8px', border: '1px solid #d9d9d9', borderRadius: 6 }}
            placeholder="Комментарий (необязательно)"
            value={comment}
            onChange={e => setComment(e.target.value)}
          />
          <Button type="primary" loading={updating} onClick={handleStatusUpdate}>
            Обновить статус
          </Button>
        </Space>
        <Divider />
      </Card>
    </div>
  );
};
