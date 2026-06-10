import type { ReactNode } from 'react';
import { useIsAuthenticated } from '@refinedev/core';
import { Spin } from 'antd';

const authLoading = (
  <div style={{ display: 'flex', minHeight: '100vh', alignItems: 'center', justifyContent: 'center' }}>
    <Spin size="large" />
  </div>
);

/** Protected shell: shows children only when check() returned authenticated. */
export function RequireLogin({
  fallback,
  children,
}: {
  fallback: ReactNode;
  children: ReactNode;
}) {
  const { data, isFetching } = useIsAuthenticated();
  if (isFetching) return authLoading;
  if (!data?.authenticated) return <>{fallback}</>;
  return <>{children}</>;
}

/** Login route: shows Login until authenticated; then renders children (e.g. redirect home). */
export function GuestRoute({
  fallback,
  children,
}: {
  fallback: ReactNode;
  children: ReactNode;
}) {
  const { data, isFetching } = useIsAuthenticated();
  if (isFetching) return authLoading;
  if (data?.authenticated) return <>{children}</>;
  return <>{fallback}</>;
}
