import { useShow } from '@refinedev/core';
import { Show } from '@refinedev/antd';
import { Typography, Space, Tag, Image, Descriptions, Switch, Divider } from 'antd';
import type { Product } from '../../types';
import { formatPrice } from '../../utils/api';
import { absoluteUploadUrl } from '../../utils/mediaUrl';
import dayjs from 'dayjs';

const { Title } = Typography;

export const ProductShow = () => {
  const { queryResult } = useShow<Product>();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={4}>{record?.name}</Title>
      
      {record?.images && record.images.length > 0 && (
        <>
          <Divider>Изображения</Divider>
          <Space size="middle" wrap>
            {record.images.map((img) => (
              <div key={img.id}>
                <Image
                  src={absoluteUploadUrl(img.image_url)}
                  alt={record.name}
                  width={200}
                  height={200}
                  style={{ objectFit: 'cover', borderRadius: 8 }}
                />
                {img.is_main && (
                  <Tag color="blue" style={{ marginTop: 8 }}>
                    Главное
                  </Tag>
                )}
              </div>
            ))}
          </Space>
        </>
      )}

      <Divider>Основная информация</Divider>
      
      <Descriptions bordered column={2}>
        <Descriptions.Item label="ID">
          <code>{record?.id}</code>
        </Descriptions.Item>
        <Descriptions.Item label="Активен">
          <Switch checked={record?.is_active} disabled />
        </Descriptions.Item>
        <Descriptions.Item label="Название" span={2}>
          {record?.name}
        </Descriptions.Item>
        <Descriptions.Item label="Описание" span={2}>
          <span style={{ whiteSpace: 'pre-line' }}>{record?.description || '-'}</span>
        </Descriptions.Item>
        <Descriptions.Item label="Цена">
          {record?.price_kopeks && formatPrice(record.price_kopeks)}
        </Descriptions.Item>
        <Descriptions.Item label="Остаток на складе">
          <Tag color={record && record.stock_quantity > 0 ? 'green' : 'red'}>
            {record?.stock_quantity} шт.
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Категория">
          {record?.category_id}
        </Descriptions.Item>
        <Descriptions.Item label="Бренд">
          {record?.firm_id || '-'}
        </Descriptions.Item>
        <Descriptions.Item label="Дата создания" span={2}>
          {record?.created_at && dayjs(record.created_at).format('DD.MM.YYYY HH:mm')}
        </Descriptions.Item>
      </Descriptions>

      {record?.tags && record.tags.length > 0 && (
        <>
          <Divider>Теги</Divider>
          <Space size="small" wrap>
            {record.tags.map((tag) => (
              <Tag key={tag.id} color={tag.color} style={{ fontSize: 14, padding: '4px 12px' }}>
                {tag.name}
              </Tag>
            ))}
          </Space>
        </>
      )}
    </Show>
  );
};
