export const catalogCategories = [
  { id: 'consoles', name: 'Игровые приставки', image: 'console.png' },
  { id: 'laptops', name: 'Ноутбуки', image: 'laptop.png' },
  { id: 'phones', name: 'Телефоны', image: 'phone.png' },
  { id: 'tablets', name: 'Планшеты', image: 'tablet.png' },
];

export const products = {
  consoles: [
    {
      id: 'xbox-series-x',
      name: 'Xbox Series X 1TB',
      title: 'Xbox Series X 1TB',
      price: 69000,
      image: 'xbox-x.png',
      rating: 5.0,
      inStock: true,
    },
    {
      id: 'xbox-series-s',
      name: 'Xbox Series S 512GB',
      title: 'Xbox Series S 512GB',
      price: 39000,
      image: 'xbox-s.png',
      rating: 5.0,
      inStock: true,
    },
  ],
};

export const subcategories = {
  consoles: [
    { id: "xbox", name: "Xbox", image: "xbox-x.png" },
    { id: "ps", name: "Sony PlayStation", image: "console.png" },
  ],
};