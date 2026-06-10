import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag, Image } from 'antd';
import { EditOutlined, PlusOutlined } from '@ant-design/icons';
import { Category } from '../../types';
import { absoluteUploadUrl } from '../../utils/mediaUrl';
import dayjs from 'dayjs';

export const CategoryList = () => {
  const { data, isLoading } = useList<Category>({
    resource: 'categories',
  });

  // parent_id у корневых в JSON часто отсутствует (omitempty), а не null — сравниваем через ??
  const sameParent = (catParent: string | null | undefined, parentId: string | null) =>
    (catParent ?? null) === parentId;

  const buildTree = (categories: Category[], parentId: string | null = null): any[] => {
    return categories
      .filter((cat) => sameParent(cat.parent_id, parentId))
      .map(cat => ({
        ...cat,
        key: cat.id,
        children: buildTree(categories, cat.id),
      }));
  };

  const treeData = data?.data ? buildTree(data.data) : [];

  const columns = [
    {
      title: 'Название',
      dataIndex: 'name',
      key: 'name',
      width: 300,
      render: (text: string, record: Category) => (
        <Space>
          {record.image_url && (
            <Image
              src={absoluteUploadUrl(record.image_url)}
              alt={text}
              width={40}
              height={40}
              style={{ objectFit: 'cover', borderRadius: 4 }}
            />
          )}
          <span style={{ fontWeight: record.parent_id ? 400 : 600 }}>{text}</span>
        </Space>
      ),
    },
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 300,
      render: (id: string) => <code style={{ fontSize: 11 }}>{id}</code>,
    },
    {
      title: 'Активна',
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
      render: (_: any, record: Category) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/categories/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Категории"
      canCreate
      createButtonProps={{
        children: 'Создать категорию',
        icon: <PlusOutlined />,
      }}
    >
      <Table
        dataSource={treeData}
        columns={columns}
        loading={isLoading}
        pagination={false}
        expandable={{
          defaultExpandAllRows: true,
        }}
      />
    </List>
  );
};
