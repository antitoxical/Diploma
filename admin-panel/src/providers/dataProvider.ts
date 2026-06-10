import type { DataProvider } from '@refinedev/core';
import { axiosInstance } from '../utils/api';
import type { ApiResponse, PaginatedResponse } from '../types';

const readResourcePath = (resource: string, id?: string | number) => {
  const normalizedResource = resource.replace(/^\/+|\/+$/g, '');

  // For list (no id): resources that have a special list URL
  const listMap: Record<string, string> = {
    products:    '/admin/products',
    categories:  '/admin/categories',
    firms:       '/admin/firms',
    tags:        '/admin/tags',
    users:       '/admin/users',
    orders:      '/admin/orders/all',
    promocodes:  '/admin/promocodes',
    'audit-logs': '/admin/audit-logs',
  };

  // For single record (with id): resources whose single-record URL differs from list URL
  const singleMap: Record<string, string> = {
    orders: '/admin/orders',
  };

  if (id !== undefined) {
    const base = singleMap[normalizedResource] ?? listMap[normalizedResource] ?? `/admin/${normalizedResource}`;
    return `${base}/${id}`;
  }

  return listMap[normalizedResource] ?? `/admin/${normalizedResource}`;
};

export const dataProvider: DataProvider = {
  getList: async ({ resource, pagination, filters, sorters, meta }) => {
    const url = readResourcePath(resource);

    const { current = 1, pageSize = 10 } = pagination ?? {};

    const params: any = {
      page: current,
      limit: pageSize,
    };

    // Add filters
    if (filters) {
      filters.forEach((filter) => {
        if (filter.operator === 'eq' && filter.value !== undefined) {
          params[filter.field] = filter.value;
        }
        if (filter.operator === 'contains' && filter.value) {
          params[filter.field] = filter.value;
        }
      });
    }

    // Add sorting
    if (sorters && sorters.length > 0) {
      const sorter = sorters[0];
      params.sort = sorter.field;
      params.order = sorter.order;
    }

    try {
      const response = await axiosInstance.get<ApiResponse<PaginatedResponse<any> | any[]>>(url, { params });
      const responseData = response.data.data;

      if (!responseData) {
        return { data: [], total: 0 };
      }

      const items = Array.isArray(responseData) ? responseData : (responseData.items ?? []);
      const total = Array.isArray(responseData) ? responseData.length : (responseData.total ?? items.length);

      return { data: items, total };
    } catch (error) {
      console.error('getList error:', error);
      throw error;
    }
  },

  getOne: async ({ resource, id }) => {
    const url = readResourcePath(resource, id);
    
    try {
      const response = await axiosInstance.get<ApiResponse<any>>(url);
      return {
        data: response.data.data,
      };
    } catch (error) {
      console.error('getOne error:', error);
      throw error;
    }
  },

  create: async ({ resource, variables }) => {
    const url = `/admin/${resource}`;
    
    try {
      const response = await axiosInstance.post<ApiResponse<any>>(url, variables);
      return {
        data: response.data.data,
      };
    } catch (error) {
      console.error('create error:', error);
      throw error;
    }
  },

  update: async ({ resource, id, variables }) => {
    const url = `/admin/${resource}/${id}`;
    
    try {
      const response = await axiosInstance.put<ApiResponse<any>>(url, variables);
      return {
        data: response.data.data,
      };
    } catch (error) {
      console.error('update error:', error);
      throw error;
    }
  },

  deleteOne: async ({ resource, id }) => {
    const url = `/admin/${resource}/${id}`;
    
    try {
      const response = await axiosInstance.delete<ApiResponse<any>>(url);
      return {
        data: response.data.data,
      };
    } catch (error) {
      console.error('deleteOne error:', error);
      throw error;
    }
  },

  getApiUrl: () => {
    return import.meta.env.VITE_API_URL || 'http://localhost:8090/api/v1';
  },

  custom: async ({ url, method, payload, query, headers }) => {
    try {
      const response = await axiosInstance({
        url,
        method,
        data: payload,
        params: query,
        headers,
      });

      const body = response.data as ApiResponse<unknown> | Record<string, unknown>;
      const inner =
        body &&
        typeof body === 'object' &&
        'status' in body &&
        'data' in body &&
        (body as ApiResponse<unknown>).data !== undefined
          ? (body as ApiResponse<unknown>).data
          : body;

      return {
        data: inner,
      };
    } catch (error) {
      console.error('custom error:', error);
      throw error;
    }
  },
};
