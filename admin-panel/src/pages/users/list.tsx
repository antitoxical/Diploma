import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Tag } from 'antd';
import { User } from '../../types';
import dayjs from 'dayjs';

export const UserList = () => {
  const { data, isLoading } = useList<User>({
    resource: 'users',
  });

  const columns = [
    {
      title: 'Telegram ID',
      dataIndex: 'telegram_id',
      key: 'telegram_id',
      render: (id: number) => <code>{id}</code>,
    },
    {
      title: 'Имя пользователя',
      dataIndex: 'username',
      key: 'username',
      render: (username: string | null) => username || '-',
    },
    {
      title: 'Имя',
      dataIndex: 'first_name',
      key: 'first_name',
    },
    {
      title: 'Фамилия',
      dataIndex: 'last_name',
      key: 'last_name',
    },
    {
      title: 'Телефон',
      dataIndex: 'phone',
      key: 'phone',
      render: (phone: string | null) => phone || '-',
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
      title: 'Дата регистрации',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm'),
    },
  ];

  return (
    <List title="Пользователи">
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
