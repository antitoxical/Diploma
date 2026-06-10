import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { Promocode } from '../../types';
import dayjs from 'dayjs';

export const PromocodeList = () => {
  const { data, isLoading } = useList<Promocode>({
    resource: 'promocodes',
  });

  const columns = [
    {
      title: 'Код',
      dataIndex: 'code',
      key: 'code',
      render: (code: string) => <code style={{ fontSize: 13, fontWeight: 600 }}>{code}</code>,
    },
    {
      title: 'Скидка',
      dataIndex: 'discount_percent',
      key: 'discount_percent',
      render: (discount: number) => `${discount}%`,
    },
    {
      title: 'Использовано / Максимум',
      key: 'uses',
      render: (_: any, record: Promocode) => (
        `${record.current_uses} / ${record.max_uses || '∞'}`
      ),
    },
    {
      title: 'Истекает',
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (date: string | null) => date ? dayjs(date).format('DD.MM.YYYY') : 'Не ограничен',
    },
    {
      title: 'Активен',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? 'Да' : 'Нет'}
        </Tag>
      ),
    },
    {
      title: 'Дата создания',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY'),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_: any, record: Promocode) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/promocodes/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Промокоды"
      canCreate
      createButtonProps={{
        children: 'Создать промокод',
      }}
    >
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
