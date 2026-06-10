import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button } from 'antd';
import { EditOutlined, PlusOutlined } from '@ant-design/icons';
import type { Firm } from '../../types';
import dayjs from 'dayjs';

export const FirmList = () => {
  const { data, isLoading } = useList<Firm>({
    resource: 'firms',
  });

  const columns = [
    {
      title: 'Название бренда',
      dataIndex: 'name',
      key: 'name',
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
      render: (_: any, record: Firm) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/firms/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Бренды"
      canCreate
      createButtonProps={{
        children: 'Создать бренд',
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
