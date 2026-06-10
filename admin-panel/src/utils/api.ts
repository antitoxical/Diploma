import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8090/api/v1';

export const axiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include JWT token
axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle errors
axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    // Do not force redirect here. Refine authProvider handles auth flow.
    return Promise.reject(error);
  }
);

// kopeks → основная единица валюты (как у BYN: 100 копеек = 1 Br)
export const kopeksToRubles = (kopeks: number): number => {
  return kopeks / 100;
};

// основная единица → копейки для API
export const rublesToKopeks = (rubles: number): number => {
  return Math.round(rubles * 100);
};

// Helper function to format price display
export const formatPrice = (kopeks: number): string => {
  const major = kopeksToRubles(kopeks);
  return `${major.toLocaleString('be-BY', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} Br`;
};
