import { Upload, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { UploadChangeParam, UploadFile, UploadProps } from 'antd/es/upload/interface';
import { useEffect, useState } from 'react';
import { axiosInstance } from '../utils/api';
import { absoluteUploadUrl } from '../utils/mediaUrl';
import type { ApiResponse, ProductImage } from '../types';

function toUrlStrings(items: unknown[]): string[] {
  return items
    .map((item) => {
      if (typeof item === 'string') return item;
      if (item && typeof item === 'object' && 'image_url' in item) {
        return String((item as ProductImage).image_url);
      }
      return '';
    })
    .filter(Boolean);
}

function extractUploadedUrl(response: unknown): string | undefined {
  if (!response || typeof response !== 'object') return undefined;
  const r = response as Record<string, unknown>;
  if (typeof r.url === 'string') return r.url;
  if (r.data && typeof r.data === 'object') {
    const d = r.data as Record<string, unknown>;
    if (typeof d.url === 'string') return d.url;
  }
  return undefined;
}

function urlsToFileList(urls: string[]): UploadFile[] {
  return urls.map((url, index) => ({
    uid: `stored-${index}-${url}`,
    name: `image-${index + 1}`,
    status: 'done' as const,
    url: absoluteUploadUrl(url),
    response: { url },
  }));
}

function fileListToUrls(list: UploadFile[]): string[] {
  return list
    .filter((f) => f.status === 'done')
    .map((f) => extractUploadedUrl(f.response))
    .filter((url): url is string => Boolean(url));
}

interface ImageUploadProps {
  value?: string[];
  onChange?: (value: string[]) => void;
  maxCount?: number;
}

export const ImageUpload = ({ value = [], onChange, maxCount = 5 }: ImageUploadProps) => {
  const urls = toUrlStrings(Array.isArray(value) ? value : []);
  const [fileList, setFileList] = useState<UploadFile[]>(() => urlsToFileList(urls));

  useEffect(() => {
    setFileList((prev) => {
      const uploading = prev.filter((f) => f.status === 'uploading');
      return [...urlsToFileList(urls), ...uploading];
    });
  }, [urls.join('|')]);

  const handleChange: UploadProps['onChange'] = (info: UploadChangeParam<UploadFile>) => {
    const next = [...info.fileList].slice(-maxCount);
    setFileList(next);

    if (info.file.status === 'done') {
      if (extractUploadedUrl(info.file.response)) {
        message.success('Изображение загружено');
      }
    }
    if (info.file.status === 'error') {
      message.error('Ошибка загрузки изображения');
    }

    if (info.file.status === 'done' || info.file.status === 'removed') {
      onChange?.(fileListToUrls(next));
    }
  };

  const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    const isImage = file.type.startsWith('image/');
    if (!isImage) {
      message.error('Можно загружать только изображения!');
      return Upload.LIST_IGNORE;
    }
    const isLt5M = file.size / 1024 / 1024 < 5;
    if (!isLt5M) {
      message.error('Изображение должно быть меньше 5MB!');
      return Upload.LIST_IGNORE;
    }
    return true;
  };

  const uploadButton = (
    <div>
      <PlusOutlined />
      <div style={{ marginTop: 8 }}>Загрузить</div>
    </div>
  );

  return (
    <Upload
      name="file"
      listType="picture-card"
      fileList={fileList}
      onChange={handleChange}
      beforeUpload={beforeUpload}
      maxCount={maxCount}
      customRequest={async ({ file, onSuccess, onError }) => {
        const fd = new FormData();
        fd.append('file', file as File);
        try {
          const res = await axiosInstance.post<ApiResponse<{ url: string }>>(
            '/admin/uploads/image',
            fd,
          );
          const url = res.data?.data?.url;
          if (!url) {
            onError?.(new Error('Пустой ответ сервера'));
            return;
          }
          onSuccess?.({ url }, file as UploadFile);
        } catch (e) {
          onError?.(e as Error);
        }
      }}
    >
      {fileList.length >= maxCount ? null : uploadButton}
    </Upload>
  );
};
