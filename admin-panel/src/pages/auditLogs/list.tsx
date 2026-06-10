import { useList } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Tag, Typography } from 'antd';
import type { AuditLog } from '../../types/auditLog';
import dayjs from 'dayjs';

const { Text } = Typography;

const actionColors: Record<string, string> = {
  create: 'green',
  update: 'blue',
  delete: 'red',
  login: 'cyan',
  logout: 'default',
};

function parseDetails(details: Record<string, unknown> | string | null): Record<string, unknown> | null {
  if (!details) return null;
  if (typeof details === 'string') {
    try {
      return JSON.parse(details) as Record<string, unknown>;
    } catch {
      return null;
    }
  }
  return details;
}

export const AuditLogList = () => {
  const { data, isLoading } = useList<AuditLog>({
    resource: 'audit-logs',
    pagination: {
      current: 1,
      pageSize: 50,
    },
  });

  const columns = [
    {
      title: 'Дата и время',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm:ss'),
      sorter: true,
      defaultSortOrder: 'descend' as const,
    },
    {
      title: 'Администратор',
      dataIndex: 'admin_username',
      key: 'admin_username',
      render: (_: string | null, record: AuditLog) => {
        if (record.admin_username) {
          return record.admin_username;
        }
        if (record.admin_id) {
          return <code style={{ fontSize: 11 }}>{record.admin_id.substring(0, 8)}...</code>;
        }
        return '-';
      },
    },
    {
      title: 'Действие',
      dataIndex: 'action',
      key: 'action',
      render: (action: string) => (
        <Tag color={actionColors[action] || 'default'}>
          {action.toUpperCase()}
        </Tag>
      ),
      filters: [
        { text: 'CREATE', value: 'create' },
        { text: 'UPDATE', value: 'update' },
        { text: 'DELETE', value: 'delete' },
        { text: 'LOGIN', value: 'login' },
        { text: 'LOGOUT', value: 'logout' },
      ],
    },
    {
      title: 'Таблица',
      dataIndex: 'table_name',
      key: 'table_name',
      render: (tableName: string | null) => (
        <Text code>{tableName || '-'}</Text>
      ),
    },
    {
      title: 'ID записи',
      dataIndex: 'record_id',
      key: 'record_id',
      render: (recordId: string | null) => (
        recordId ? <code style={{ fontSize: 11 }}>{recordId.substring(0, 8)}...</code> : '-'
      ),
    },
    {
      title: 'Детали',
      dataIndex: 'details',
      key: 'details',
      render: (details: Record<string, unknown> | string | null) => {
        const parsed = parseDetails(details);
        if (!parsed || Object.keys(parsed).length === 0) return '-';

        return (
          <div style={{ maxWidth: 400, maxHeight: 100, overflow: 'auto' }}>
            <pre style={{ fontSize: 11, margin: 0 }}>
              {JSON.stringify(parsed, null, 2)}
            </pre>
          </div>
        );
      },
    },
  ];

  return (
    <List title="Журнал действий администраторов">
      <Table
        dataSource={data?.data}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={{
          total: data?.total,
          pageSize: 50,
          showSizeChanger: true,
          showTotal: (total) => `Всего записей: ${total}`,
        }}
        scroll={{ x: 1200 }}
      />
    </List>
  );
};
