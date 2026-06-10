import { Refine } from '@refinedev/core';
import { RefineKbar, RefineKbarProvider } from '@refinedev/kbar';
import {
  ErrorComponent,
  ThemedLayout,
  ThemedTitle,
} from '@refinedev/antd';
import routerBindings from '@refinedev/react-router-v6';
import { CatchAllNavigate } from '@refinedev/react-router-v6';
import { DocumentTitleHandler } from '@refinedev/react-router-v6';
import { UnsavedChangesNotifier } from '@refinedev/react-router-v6';
import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import { App as AntdApp, ConfigProvider } from 'antd';
import { useNotificationProvider } from '@refinedev/antd';
import ruRU from 'antd/locale/ru_RU';
import 'antd/dist/reset.css';

import { authProvider } from './providers/authProvider';
import { dataProvider } from './providers/dataProvider';
import { i18nProvider } from './providers/i18nProvider';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { ProductList, ProductCreate, ProductEdit, ProductShow } from './pages/products';
import { CategoryList, CategoryCreate, CategoryEdit } from './pages/categories';
import { OrderList, OrderShow } from './pages/orders';
import { UserList } from './pages/users';
import { PromocodeList, PromocodeCreate, PromocodeEdit } from './pages/promocodes';
import { TagList, TagCreate, TagEdit } from './pages/tags';
import { ReviewList } from './pages/reviews';
import { FirmList, FirmCreate, FirmEdit } from './pages/firms';
import { AuditLogList } from './pages/auditLogs';
import { AdminList, AdminCreate, AdminEdit } from './pages/admins';
import { DeliveryMethodList, DeliveryMethodCreate, DeliveryMethodEdit } from './pages/deliveryMethods';
import { PaymentMethodList, PaymentMethodCreate, PaymentMethodEdit } from './pages/paymentMethods';
import { GuestRoute, RequireLogin } from './components/authGuards';
import {
  ShoppingOutlined,
  AppstoreOutlined,
  TagsOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  GiftOutlined,
  StarOutlined,
  SettingOutlined,
} from '@ant-design/icons';

function App() {
  return (
    <BrowserRouter>
      <ConfigProvider locale={ruRU}>
        <AntdApp>
          <RefineKbarProvider>
            <Refine
              dataProvider={dataProvider}
              authProvider={authProvider}
              i18nProvider={i18nProvider}
              routerProvider={routerBindings}
              notificationProvider={useNotificationProvider}
              resources={[
                {
                  name: 'dashboard',
                  list: '/',
                  meta: {
                    label: 'Главная',
                    icon: <AppstoreOutlined />,
                  },
                },
                {
                  name: 'products',
                  list: '/products',
                  create: '/products/create',
                  edit: '/products/edit/:id',
                  show: '/products/show/:id',
                  meta: {
                    label: 'Товары',
                    icon: <ShoppingOutlined />,
                  },
                },
                {
                  name: 'categories',
                  list: '/categories',
                  create: '/categories/create',
                  edit: '/categories/edit/:id',
                  meta: {
                    label: 'Категории',
                    icon: <AppstoreOutlined />,
                  },
                },
                {
                  name: 'firms',
                  list: '/firms',
                  create: '/firms/create',
                  edit: '/firms/edit/:id',
                  meta: {
                    label: 'Бренды',
                    icon: <TagsOutlined />,
                  },
                },
                {
                  name: 'tags',
                  list: '/tags',
                  create: '/tags/create',
                  edit: '/tags/edit/:id',
                  meta: {
                    label: 'Теги',
                    icon: <TagsOutlined />,
                  },
                },
                {
                  name: 'orders',
                  list: '/orders',
                  show: '/orders/show/:id',
                  meta: {
                    label: 'Заказы',
                    icon: <ShoppingCartOutlined />,
                  },
                },
                {
                  name: 'users',
                  list: '/users',
                  show: '/users/show/:id',
                  meta: {
                    label: 'Пользователи',
                    icon: <UserOutlined />,
                  },
                },
                {
                  name: 'promocodes',
                  list: '/promocodes',
                  create: '/promocodes/create',
                  edit: '/promocodes/edit/:id',
                  meta: {
                    label: 'Промокоды',
                    icon: <GiftOutlined />,
                  },
                },
                {
                  name: 'reviews',
                  list: '/reviews',
                  meta: {
                    label: 'Отзывы',
                    icon: <StarOutlined />,
                  },
                },
                {
                  name: 'delivery-methods',
                  list: '/delivery-methods',
                  create: '/delivery-methods/create',
                  edit: '/delivery-methods/edit/:id',
                  meta: {
                    label: 'Доставка',
                    icon: <SettingOutlined />,
                  },
                },
                {
                  name: 'payment-methods',
                  list: '/payment-methods',
                  create: '/payment-methods/create',
                  edit: '/payment-methods/edit/:id',
                  meta: {
                    label: 'Оплата',
                    icon: <SettingOutlined />,
                  },
                },
                {
                  name: 'audit-logs',
                  list: '/audit-logs',
                  meta: {
                    label: 'Аудит-логи',
                    icon: <SettingOutlined />,
                  },
                },
                {
                  name: 'admins',
                  list: '/admins',
                  create: '/admins/create',
                  edit: '/admins/edit/:id',
                  meta: {
                    label: 'Администраторы',
                    icon: <UserOutlined />,
                  },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                useNewQueryKeys: true,
              }}
            >
              <Routes>
                <Route
                  element={
                    <RequireLogin
                      key="authenticated-inner"
                      fallback={<CatchAllNavigate to="/login" />}
                    >
                      <ThemedLayout
                        Title={({ collapsed }) => (
                          <ThemedTitle
                            collapsed={collapsed}
                            text="AlphaSales Admin"
                          />
                        )}
                      >
                        <Outlet />
                      </ThemedLayout>
                    </RequireLogin>
                  }
                >
                  <Route index element={<Dashboard />} />
                  <Route path="/products">
                    <Route index element={<ProductList />} />
                    <Route path="create" element={<ProductCreate />} />
                    <Route path="edit/:id" element={<ProductEdit />} />
                    <Route path="show/:id" element={<ProductShow />} />
                  </Route>
                  <Route path="/categories">
                    <Route index element={<CategoryList />} />
                    <Route path="create" element={<CategoryCreate />} />
                    <Route path="edit/:id" element={<CategoryEdit />} />
                  </Route>
                  <Route path="/firms">
                    <Route index element={<FirmList />} />
                    <Route path="create" element={<FirmCreate />} />
                    <Route path="edit/:id" element={<FirmEdit />} />
                  </Route>
                  <Route path="/orders">
                    <Route index element={<OrderList />} />
                    <Route path="show/:id" element={<OrderShow />} />
                  </Route>
                  <Route path="/users">
                    <Route index element={<UserList />} />
                  </Route>
                  <Route path="/promocodes">
                    <Route index element={<PromocodeList />} />
                    <Route path="create" element={<PromocodeCreate />} />
                    <Route path="edit/:id" element={<PromocodeEdit />} />
                  </Route>
                  <Route path="/tags">
                    <Route index element={<TagList />} />
                    <Route path="create" element={<TagCreate />} />
                    <Route path="edit/:id" element={<TagEdit />} />
                  </Route>
                  <Route path="/reviews">
                    <Route index element={<ReviewList />} />
                  </Route>
                  <Route path="/delivery-methods">
                    <Route index element={<DeliveryMethodList />} />
                    <Route path="create" element={<DeliveryMethodCreate />} />
                    <Route path="edit/:id" element={<DeliveryMethodEdit />} />
                  </Route>
                  <Route path="/payment-methods">
                    <Route index element={<PaymentMethodList />} />
                    <Route path="create" element={<PaymentMethodCreate />} />
                    <Route path="edit/:id" element={<PaymentMethodEdit />} />
                  </Route>
                  <Route path="/audit-logs">
                    <Route index element={<AuditLogList />} />
                  </Route>
                  <Route path="/admins">
                    <Route index element={<AdminList />} />
                    <Route path="create" element={<AdminCreate />} />
                    <Route path="edit/:id" element={<AdminEdit />} />
                  </Route>
                  <Route path="*" element={<ErrorComponent />} />
                </Route>
                <Route
                  path="/login"
                  element={
                    <GuestRoute key="authenticated-outer" fallback={<Login />}>
                      <CatchAllNavigate to="/" />
                    </GuestRoute>
                  }
                />
              </Routes>
              <RefineKbar />
              <UnsavedChangesNotifier />
              <DocumentTitleHandler />
            </Refine>
          </RefineKbarProvider>
        </AntdApp>
      </ConfigProvider>
    </BrowserRouter>
  );
}

export default App;
