import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { BrowserRouter } from 'react-router-dom';
import { TelegramProvider } from './context/TelegramContext';
import { ThemeProvider } from './context/ThemeContext';
import { CartProvider } from './context/CartContext';
import { FavoritesProvider } from './context/FavoritesContext';
import { ToastProvider } from './context/ToastContext';

const container = document.getElementById('root') as HTMLElement;
const root = ReactDOM.createRoot(container);

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <TelegramProvider>
        <ThemeProvider>
          <CartProvider>
            <FavoritesProvider>
              <ToastProvider>
                <App />
              </ToastProvider>
            </FavoritesProvider>
          </CartProvider>
        </ThemeProvider>
      </TelegramProvider>
    </BrowserRouter>
  </React.StrictMode>
);
