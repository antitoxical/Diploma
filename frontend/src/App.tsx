import './App.css';
import './style/toast.css';
import { Routes, Route } from "react-router-dom";
import Home from "./pages/home";
import BottomNav from './components/BottomNav';
import Catalog from './pages/catalog';
import { CategoryPage } from './pages/category';
import Cart from './pages/basket';
import Favorites from "./pages/favorites";
import ProductPage from "./pages/product";
import Reviews from "./pages/reviews";
import Checkout from "./pages/checkout";
import Profile from "./pages/profile";
import Orders from "./pages/orders";
import OrderDetail from "./pages/orderDetail";
import Notifications from "./pages/notifications";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/catalog" element={<Catalog />} />
        <Route path="/category/:categoryId" element={<CategoryPage />} />
        <Route path="/basket" element={<Cart />} />
        <Route path="/favorites" element={<Favorites />} />
        <Route path="/product/:productId" element={<ProductPage />} />
        <Route path="/product/:productId/reviews" element={<Reviews />} />
        <Route path="/checkout" element={<Checkout />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/orders" element={<Orders />} />
        <Route path="/order/:orderId" element={<OrderDetail />} />
        <Route path="/notifications" element={<Notifications />} />
      </Routes>
      <BottomNav />
    </>
  );
}

export default App;
