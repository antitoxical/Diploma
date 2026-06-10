import type { AuthProvider } from '@refinedev/core';
import { axiosInstance } from '../utils/api';
import type { ApiResponse, LoginRequest, LoginResponse } from '../types';

export const authProvider: AuthProvider = {
  login: async ({ username, password }: LoginRequest) => {
    try {
      const response = await axiosInstance.post<ApiResponse<LoginResponse>>(
        '/admin/login',
        { username, password }
      );

      if (response.data?.data?.token) {
        localStorage.setItem('auth_token', response.data.data.token);
        localStorage.setItem('admin_user', JSON.stringify(response.data.data.admin));
        return {
          success: true,
          redirectTo: '/',
        };
      }

      return {
        success: false,
        error: {
          name: 'LoginError',
          message: 'Invalid credentials',
        },
      };
    } catch (error: any) {
      return {
        success: false,
        error: {
          name: 'LoginError',
          message: error.response?.data?.data || 'Login failed',
        },
      };
    }
  },

  logout: async () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('admin_user');
    return {
      success: true,
      redirectTo: '/login',
    };
  },

  check: async () => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      return {
        authenticated: true,
      };
    }

    return {
      authenticated: false,
      redirectTo: '/login',
    };
  },

  getPermissions: async () => null,

  getIdentity: async () => {
    const adminUser = localStorage.getItem('admin_user');
    if (adminUser) {
      return JSON.parse(adminUser);
    }
    return null;
  },

  onError: async (error) => {
    if (error.response?.status === 401) {
      return {
        logout: true,
        redirectTo: '/login',
        error,
      };
    }

    return { error };
  },
};
