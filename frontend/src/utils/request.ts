export async function api(path: string, options: RequestInit = {}) {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    'Cache-Control': 'no-store',
    Pragma: 'no-cache',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(options.headers || {}),
  } as any;
  const res = await fetch(`/api/v1${path}`, { ...options, headers, cache: 'no-store' });
  let data: any = {};
  try { data = await res.json(); } catch {}
  if (res.status === 401 || data?.code === 40101) {
    // 令牌失效或被替换（例如后端更换密钥），清除后回到登录
    localStorage.removeItem('token');
    try { window.dispatchEvent(new CustomEvent('auth-logout', { detail: { reason: 'unauthorized' } })); } catch {}
    throw new Error('未授权或登录失效，请重新登录');
  }
  if (!res.ok || data.code !== 0) {
    throw new Error(data?.message || `请求失败: ${res.status}`);
  }
  return data.data;
}
