import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { Order } from '../../types';
import { formatPrice } from '../../utils/api';
import dayjs from 'dayjs';

const statusColors: Record<string, string> = {
  created: 'blue',
  paid: 'green',
  delivering: 'orange',
  completed: 'success',
  cancelled: 'red',
};

const statusLabels: Record<string, string> = {
  created: 'Создан',
  paid: 'Оплачен',
  delivering: 'Доставляется',
  completed: 'Завершён',
  cancelled: 'Отменён',
};

export const OrderList = () => {
  const { data, isLoading } = useList<Order>({
    resource: 'orders',
  });

  const columns = [
    {
      title: 'ID заказа',
      dataIndex: 'id',
      key: 'id',
      width: 100,
      render: (id: string) => <code style={{ fontSize: 11 }}>{id.substring(0, 8)}...</code>,
    },
    {
      title: 'Пользователь',
      dataIndex: 'user_id',
      key: 'user_id',
      render: (userId: string) => <code style={{ fontSize: 11 }}>{userId.substring(0, 8)}...</code>,
    },
    {
      title: 'Статус',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status]}>
          {statusLabels[status] || status}
        </Tag>
      ),
      filters: [
        { text: 'Создан', value: 'created' },
        { text: 'Оплачен', value: 'paid' },
        { text: 'Доставляется', value: 'delivering' },
        { text: 'Завершён', value: 'completed' },
        { text: 'Отменён', value: 'cancelled' },
      ],
    },
    {
      title: 'Сумма',
      dataIndex: 'total_amount_kopeks',
      key: 'total_amount_kopeks',
      render: (amount: number) => formatPrice(amount),
      sorter: true,
    },
    {
      title: 'Адрес доставки',
      dataIndex: 'address_snapshot',
      key: 'address_snapshot',
      ellipsis: true,
      width: 200,
    },
    {
      title: 'Дата создания',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm'),
      sorter: true,
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_: any, record: Order) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            href={`/orders/show/${record.id}`}
          >
            Просмотр
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List title="Заказы">
      <Table
        dataSource={data?.data}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={{
          total: data?.total,
          pageSize: 10,
          showSizeChanger: true,
          showTotal: (total) => `Всего: ${total}`,
        }}
      />
    </List>
  );
};
