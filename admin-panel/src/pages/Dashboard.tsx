import { Card, Col, Row, Statistic, Typography, Spin } from 'antd';
import { ShoppingCartOutlined, UserOutlined, DollarOutlined, TagsOutlined } from '@ant-design/icons';
import { useCustom } from '@refinedev/core';
import { LineChart, Line, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, Legend } from 'recharts';
import { formatPrice } from '../utils/api';
import dayjs from 'dayjs';

const { Title } = Typography;

interface DashboardStats {
  total_orders: number;
  total_users: number;
  total_revenue_kopeks: number;
  total_products: number;
  orders_by_day: Array<{ date: string; count: number; revenue_kopeks: number }>;
  orders_by_status: Array<{ status: string; count: number }>;
  top_products: Array<{ id: string; name: string; sold_count: number; revenue_kopeks: number }>;
}

export const Dashboard = () => {
  const { data, isFetching } = useCustom<DashboardStats>({
    url: '/admin/dashboard/stats',
    method: 'get',
    queryOptions: {
      retry: false,
      enabled: true,
    },
  });

  const stats = data?.data;

  if (isFetching) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
        <div style={{ marginTop: 16 }}>Загрузка статистики...</div>
      </div>
    );
  }

  const ordersByDay = stats?.orders_by_day?.map(item => ({
    ...item,
    date: dayjs(item.date).format('DD.MM'),
    revenue: item.revenue_kopeks / 100,
  })) || [];

  const statusLabels: Record<string, string> = {
    created: 'Создан',
    paid: 'Оплачен',
    delivering: 'Доставляется',
    completed: 'Завершён',
    cancelled: 'Отменён',
  };

  const ordersByStatus = stats?.orders_by_status?.map(item => ({
    status: statusLabels[item.status] || item.status,
    count: item.count,
  })) || [];

  return (
    <div>
      <Title level={2}>Панель управления</Title>
      
      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Всего заказов"
              value={stats?.total_orders || 0}
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Пользователей"
              value={stats?.total_users || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Общая выручка"
              value={stats?.total_revenue_kopeks ? (stats.total_revenue_kopeks / 100).toFixed(2) : 0}
              suffix="Br"
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Товаров"
              value={stats?.total_products || 0}
              prefix={<TagsOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={16}>
          <Card title="Заказы и выручка за последние 30 дней">
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={ordersByDay}>
                <defs>
                  <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#1890ff" stopOpacity={0.8}/>
                    <stop offset="95%" stopColor="#1890ff" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Legend />
                <Area
                  yAxisId="left"
                  type="monotone"
                  dataKey="count"
                  stroke="#52c41a"
                  fillOpacity={1}
                  fill="url(#colorRevenue)"
                  name="Заказы"
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="revenue"
                  stroke="#ff4d4f"
                  strokeWidth={2}
                  name="Выручка (Br)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card title="Распределение заказов по статусам">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={ordersByStatus}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="status" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="count" fill="#1890ff" name="Количество" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      {stats?.top_products && stats.top_products.length > 0 && (
        <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
          <Col span={24}>
            <Card title="Топ товаров по продажам">
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={stats.top_products} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis type="category" dataKey="name" width={200} />
                  <Tooltip formatter={(value: any, name: string) => {
                    if (name === 'Выручка (Br)') {
                      return `${(value as number).toFixed(2)} Br`;
                    }
                    return value;
                  }} />
                  <Legend />
                  <Bar dataKey="sold_count" fill="#52c41a" name="Продано (шт.)" />
                  <Bar
                    dataKey="revenue_kopeks"
                    fill="#1890ff"
                    name="Выручка (Br)"
                  />
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );
};
