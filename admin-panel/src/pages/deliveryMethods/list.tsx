import { useTable } from '@refinedev/antd';
import { List } from '@refinedev/antd';
import { Table, Tag, Space, Button } from 'antd';
import { EditOutlined, PlusOutlined } from '@ant-design/icons';
import type { DeliveryMethod } from '../../types';
import { DELIVERY_KIND_LABELS } from '../../types';
import { formatPrice } from '../../utils/api';
import { useNavigate } from 'react-router-dom';

export const DeliveryMethodList = () => {
  const navigate = useNavigate();
  const { tableProps } = useTable<DeliveryMethod>({
    resource: 'delivery-methods',
    syncWithLocation: true,
  });

  return (
    <List
      title="Способы доставки"
      canCreate
      createButtonProps={{
        children: 'Добавить',
        icon: <PlusOutlined />,
        onClick: () => navigate('/delivery-methods/create'),
      }}
    >
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="name" title="Название" />
        <Table.Column
          dataIndex="kind"
          title="Тип"
          render={(v: string) => DELIVERY_KIND_LABELS[v as keyof typeof DELIVERY_KIND_LABELS] ?? v}
        />
        <Table.Column
          dataIndex="price_kopeks"
          title="Цена"
          render={(v: number) => formatPrice(v)}
        />
        <Table.Column
          dataIndex="description"
          title="Описание"
          render={(v: string | null) => v ?? '—'}
        />
        <Table.Column
          dataIndex="is_active"
          title="Активен"
          render={(v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? 'Да' : 'Нет'}</Tag>}
        />
        <Table.Column
          title="Действия"
          render={(_: unknown, record: DeliveryMethod) => (
            <Space>
              <Button
                type="link"
                icon={<EditOutlined />}
                onClick={() => navigate(`/delivery-methods/edit/${record.id}`)}
              >
                Изменить
              </Button>
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
