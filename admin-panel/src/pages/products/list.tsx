import { useList, useMany } from '@refinedev/core';
import { List } from '@refinedev/antd';
import { Table, Space, Button, Tag, Switch, Input } from 'antd';
import { EditOutlined, EyeOutlined, DownloadOutlined, FileExcelOutlined, SearchOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { FilterDropdownProps } from 'antd/es/table/interface';
import type { Product, Category, Firm } from '../../types';
import { formatPrice, kopeksToRubles } from '../../utils/api';
import { absoluteUploadUrl } from '../../utils/mediaUrl';
import { exportToExcel, exportToCSV, flattenObject } from '../../utils/export';
import dayjs from 'dayjs';

const NameFilterDropdown = ({
  setSelectedKeys,
  selectedKeys,
  confirm,
  clearFilters,
}: FilterDropdownProps) => (
  <div style={{ padding: 8 }} onKeyDown={(e) => e.stopPropagation()}>
    <Input
      placeholder="Поиск по названию"
      value={selectedKeys[0] as string}
      onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
      onPressEnter={() => confirm()}
      style={{ marginBottom: 8, display: 'block' }}
    />
    <Space>
      <Button type="primary" onClick={() => confirm()} size="small">
        OK
      </Button>
      <Button
        onClick={() => {
          clearFilters?.();
          confirm();
        }}
        size="small"
      >
        Сброс
      </Button>
    </Space>
  </div>
);

export const ProductList = () => {
  const { data, isLoading } = useList<Product>({
    resource: 'products',
    pagination: { mode: 'off' },
  });

  const { data: categoriesList } = useList<Category>({
    resource: 'categories',
    pagination: { mode: 'off' },
  });

  const { data: firmsList } = useList<Firm>({
    resource: 'firms',
    pagination: { mode: 'off' },
  });

  const categoryIds = [
    ...new Set(data?.data?.map((item) => item.category_id).filter(Boolean) ?? []),
  ];
  const { data: categoriesData } = useMany<Category>({
    resource: 'categories',
    ids: categoryIds,
    queryOptions: {
      enabled: categoryIds.length > 0,
    },
  });

  const firmIds = [
    ...new Set(
      data?.data?.filter((item) => item.firm_id).map((item) => item.firm_id!) ?? [],
    ),
  ];
  const { data: firmsData } = useMany<Firm>({
    resource: 'firms',
    ids: firmIds,
    queryOptions: {
      enabled: firmIds.length > 0,
    },
  });

  const categoryName = (categoryId: string) =>
    categoriesData?.data?.find((cat) => cat.id === categoryId)?.name ||
    categoriesList?.data?.find((cat) => cat.id === categoryId)?.name ||
    '-';

  const firmName = (firmId: string | null) => {
    if (!firmId) return '-';
    return (
      firmsData?.data?.find((f) => f.id === firmId)?.name ||
      firmsList?.data?.find((f) => f.id === firmId)?.name ||
      '-'
    );
  };

  const handleExportExcel = () => {
    const exportData = data?.data?.map(product => {
      return flattenObject({
        id: product.id,
        name: product.name,
        description: product.description || '',
        category: categoryName(product.category_id),
        firm: firmName(product.firm_id),
        price_rubles: kopeksToRubles(product.price_kopeks),
        stock_quantity: product.stock_quantity,
        is_active: product.is_active ? 'Да' : 'Нет',
        tags: product.tags?.map(t => t.name).join(', ') || '',
        created_at: dayjs(product.created_at).format('DD.MM.YYYY HH:mm'),
      });
    }) || [];

    exportToExcel(exportData, `products_${dayjs().format('YYYY-MM-DD')}`);
  };

  const handleExportCSV = () => {
    const exportData = data?.data?.map(product => {
      return flattenObject({
        id: product.id,
        name: product.name,
        description: product.description || '',
        category: categoryName(product.category_id),
        firm: firmName(product.firm_id),
        price_rubles: kopeksToRubles(product.price_kopeks),
        stock_quantity: product.stock_quantity,
        is_active: product.is_active ? 'Да' : 'Нет',
        tags: product.tags?.map(t => t.name).join(', ') || '',
        created_at: dayjs(product.created_at).format('DD.MM.YYYY HH:mm'),
      });
    }) || [];

    exportToCSV(exportData, `products_${dayjs().format('YYYY-MM-DD')}`);
  };

  const categoryFilters =
    categoriesList?.data?.map((c) => ({ text: c.name, value: c.id })) ?? [];

  const firmFilters =
    firmsList?.data?.map((f) => ({ text: f.name, value: f.id })) ?? [];

  const columns: ColumnsType<Product> = [
    {
      title: 'Название',
      dataIndex: 'name',
      key: 'name',
      width: 300,
      filterDropdown: (props) => <NameFilterDropdown {...props} />,
      filterIcon: (filtered) => (
        <SearchOutlined style={{ color: filtered ? '#1677ff' : undefined }} />
      ),
      onFilter: (value, record) =>
        record.name.toLowerCase().includes(String(value).toLowerCase()),
      render: (text: string, record: Product) => (
        <div>
          <div style={{ fontWeight: 500 }}>{text}</div>
          {record.images && record.images.length > 0 && (
            <img
              src={absoluteUploadUrl(
                record.images.find(img => img.is_main)?.image_url || record.images[0].image_url,
              )}
              alt={text}
              style={{ width: 50, height: 50, objectFit: 'cover', marginTop: 4, borderRadius: 4 }}
            />
          )}
        </div>
      ),
    },
    {
      title: 'Категория',
      dataIndex: 'category_id',
      key: 'category_id',
      filters: categoryFilters,
      onFilter: (value, record) => record.category_id === value,
      render: (categoryId: string) => categoryName(categoryId),
    },
    {
      title: 'Бренд',
      dataIndex: 'firm_id',
      key: 'firm_id',
      filters: firmFilters,
      onFilter: (value, record) => record.firm_id === value,
      render: (firmId: string | null) => firmName(firmId),
    },
    {
      title: 'Цена',
      dataIndex: 'price_kopeks',
      key: 'price_kopeks',
      render: (price: number) => formatPrice(price),
      sorter: (a, b) => a.price_kopeks - b.price_kopeks,
    },
    {
      title: 'Остаток',
      dataIndex: 'stock_quantity',
      key: 'stock_quantity',
      render: (stock: number) => (
        <Tag color={stock > 0 ? 'green' : 'red'}>
          {stock} шт.
        </Tag>
      ),
      sorter: (a, b) => a.stock_quantity - b.stock_quantity,
      filters: [
        { text: 'В наличии', value: 'in_stock' },
        { text: 'Нет на складе', value: 'out_of_stock' },
      ],
      onFilter: (value, record) =>
        value === 'in_stock' ? record.stock_quantity > 0 : record.stock_quantity <= 0,
    },
    {
      title: 'Активен',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => <Switch checked={isActive} disabled />,
      filters: [
        { text: 'Да', value: true },
        { text: 'Нет', value: false },
      ],
      onFilter: (value, record) => record.is_active === value,
    },
    {
      title: 'Теги',
      dataIndex: 'tags',
      key: 'tags',
      render: (tags: Product['tags']) => (
        <>
          {tags?.map((tag) =>
            tag ? (
              <Tag key={tag.id} color={tag.color}>
                {tag.name}
              </Tag>
            ) : null,
          )}
        </>
      ),
    },
    {
      title: 'Дата создания',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD.MM.YYYY HH:mm'),
      sorter: (a, b) => dayjs(a.created_at).valueOf() - dayjs(b.created_at).valueOf(),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_: unknown, record: Product) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            href={`/products/show/${record.id}`}
          >
            Просмотр
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            href={`/products/edit/${record.id}`}
          >
            Изменить
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <List
      title="Товары"
      canCreate
      createButtonProps={{
        children: 'Создать товар',
      }}
      headerButtons={({ defaultButtons }) => (
        <>
          {defaultButtons}
          <Button
            type="default"
            icon={<FileExcelOutlined />}
            onClick={handleExportExcel}
          >
            Excel
          </Button>
          <Button
            type="default"
            icon={<DownloadOutlined />}
            onClick={handleExportCSV}
          >
            CSV
          </Button>
        </>
      )}
    >
      <Table
        dataSource={data?.data}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showTotal: (total) => `Всего: ${total}`,
        }}
      />
    </List>
  );
};
