import { useTable } from '@refinedev/antd';
import { List } from '@refinedev/antd';
import { Table, Tag, Space, Button } from 'antd';
import { EditOutlined, PlusOutlined } from '@ant-design/icons';
import type { PaymentMethod } from '../../types';
import { useNavigate } from 'react-router-dom';

export const PaymentMethodList = () => {
  const navigate = useNavigate();
  const { tableProps } = useTable<PaymentMethod>({
    resource: 'payment-methods',
    syncWithLocation: true,
  });

  return (
    <List
      title="Способы оплаты"
      canCreate
      createButtonProps={{
        children: 'Добавить',
        icon: <PlusOutlined />,
        onClick: () => navigate('/payment-methods/create'),
      }}
    >
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="name" title="Название" />
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
          render={(_: unknown, record: PaymentMethod) => (
            <Space>
              <Button
                type="link"
                icon={<EditOutlined />}
                onClick={() => navigate(`/payment-methods/edit/${record.id}`)}
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
