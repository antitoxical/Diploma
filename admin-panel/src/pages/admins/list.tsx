import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag } from 'antd';
import { EditOutlined, PlusOutlined } from '@ant-design/icons';
import type { Admin } from '../../types';
import dayjs from 'dayjs';

export const AdminList = () => {
  const { data, isLoading } = useList<Admin>({
    resource: 'admins',
  });

  const columns = [
    {
      title: 'Имя пользователя',
      dataIndex: 'username',
      key: 'username',
      render: (username: string) => <strong>{username}</strong>,
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
      filters: [
        { text: 'Активные', value: true },
        { text: 'Неактивные', value: false },
      ],
    },
    {
      title: 'Дата создания',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm'),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_: any, record: Admin) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/admins/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Администраторы"
      canCreate
      createButtonProps={{
        children: 'Создать администратора',
        icon: <PlusOutlined />,
      }}
    >
      <Table
        dataSource={data?.data}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={false}
      />
    </List>
  );
};
