import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Rate, Tag, message } from 'antd';
import { CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { Review } from '../../types';
import { axiosInstance } from '../../utils/api';
import dayjs from 'dayjs';

export const ReviewList = () => {
  const { data, isLoading, refetch } = useList<Review>({
    resource: 'reviews',
  });

  const handleModerate = async (reviewId: string, approved: boolean) => {
    try {
      await axiosInstance.patch(`/admin/reviews/${reviewId}/moderate`, { approved });
      message.success(approved ? 'Отзыв одобрен' : 'Отзыв отклонён');
      refetch();
    } catch {
      message.error('Ошибка при модерации отзыва');
    }
  };

  const columns = [
    {
      title: 'Пользователь',
      dataIndex: 'user_id',
      key: 'user_id',
      render: (userId: string) => <code style={{ fontSize: 11 }}>{userId.substring(0, 8)}...</code>,
    },
    {
      title: 'Товар',
      dataIndex: 'product_id',
      key: 'product_id',
      render: (productId: string) => <code style={{ fontSize: 11 }}>{productId.substring(0, 8)}...</code>,
    },
    {
      title: 'Рейтинг',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => <Rate disabled value={rating} />,
    },
    {
      title: 'Комментарий',
      dataIndex: 'comment',
      key: 'comment',
      ellipsis: true,
      width: 300,
      render: (comment: string | null) => comment || <span style={{ color: '#999' }}>Без комментария</span>,
    },
    {
      title: 'Статус',
      dataIndex: 'is_moderated',
      key: 'is_moderated',
      render: (isModerated: boolean) => (
        <Tag color={isModerated ? 'green' : 'orange'}>
          {isModerated ? 'Одобрен' : 'На модерации'}
        </Tag>
      ),
      filters: [
        { text: 'Одобренные', value: true },
        { text: 'На модерации', value: false },
      ],
    },
    {
      title: 'Дата',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm'),
      sorter: true,
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_: any, record: Review) => (
        <Space>
          {!record.is_moderated && (
            <>
              <Button
                type="link"
                icon={<CheckOutlined />}
                onClick={() => handleModerate(record.id, true)}
                style={{ color: '#52c41a' }}
              >
                Одобрить
              </Button>
              <Button
                type="link"
                danger
                icon={<CloseOutlined />}
                onClick={() => handleModerate(record.id, false)}
              >
                Удалить
              </Button>
            </>
          )}
          {record.is_moderated && (
            <Tag color="green">Обработан</Tag>
          )}
        </Space>
      ),
    },
  ];

  return (
    <List title="Отзывы">
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
