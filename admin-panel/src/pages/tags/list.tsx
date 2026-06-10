import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { Tag as TagType } from '../../types';
import dayjs from 'dayjs';

export const TagList = () => {
  const { data, isLoading } = useList<TagType>({
    resource: 'tags',
  });

  const columns = [
    {
      title: 'Название',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: TagType) => (
        <Tag color={record.color} style={{ fontSize: 14, padding: '4px 12px' }}>
          {name}
        </Tag>
      ),
    },
    {
      title: 'Цвет',
      dataIndex: 'color',
      key: 'color',
      render: (color: string) => (
        <Space>
          <div
            style={{
              width: 20,
              height: 20,
              backgroundColor: color,
              borderRadius: 4,
              border: '1px solid #ddd',
            }}
          />
          <code>{color}</code>
        </Space>
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
      render: (_: any, record: TagType) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/tags/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Теги"
      canCreate
      createButtonProps={{
        children: 'Создать тег',
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
