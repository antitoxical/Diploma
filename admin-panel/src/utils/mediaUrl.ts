/** Относительный путь API (/api/v1/...) → абсолютный URL для <img> на другом порту (Vite :5173). */
export function absoluteUploadUrl(pathOrUrl: string): string {
  if (!pathOrUrl) return '';
  if (pathOrUrl.startsWith('http://') || pathOrUrl.startsWith('https://')) return pathOrUrl;
  const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8090/api/v1';
  const origin = apiBase.replace(/\/api\/v1\/?$/, '');
  return pathOrUrl.startsWith('/') ? `${origin}${pathOrUrl}` : `${origin}/${pathOrUrl}`;
}
