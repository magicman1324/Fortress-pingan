import type {
  Asset, Session, AuditPage, DashboardStats, LoginResponse,
} from './types';

const BASE = '/api/v1';

function getToken(): string | null {
  return localStorage.getItem('token');
}

async function request<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(opts.headers as Record<string, string> || {}),
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  if (!(opts.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  const res = await fetch(`${BASE}${path}`, { ...opts, headers });

  if (res.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || `HTTP ${res.status}`);
  }
  return data as T;
}

export const api = {
  // Auth
  login: (username: string, password: string) =>
    request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  // Assets
  getAssets: (search = '') =>
    request<Asset[]>(`/assets?search=${encodeURIComponent(search)}`),

  createAsset: (data: Partial<Asset> & { name: string; host: string; username: string; credential: string }) =>
    request<Asset>('/assets', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateAsset: (id: number, data: Partial<Asset>) =>
    request<Asset>(`/assets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteAsset: (id: number) =>
    request<void>(`/assets/${id}`, { method: 'DELETE' }),

  testAsset: (id: number) =>
    request<{ status: string; error?: string }>(`/assets/${id}/test`, { method: 'POST' }),

  // Sessions
  getSessions: (page = 1, size = 20) =>
    request<Session[]>(`/sessions?page=${page}&size=${size}`),

  terminateSession: (id: number) =>
    request<void>(`/sessions/${id}/terminate`, { method: 'POST' }),

  // Audit Logs
  getAuditLogs: (params: { page?: number; size?: number; user_id?: number; asset_id?: number; keyword?: string } = {}) => {
    const qs = new URLSearchParams();
    if (params.page) qs.set('page', String(params.page));
    if (params.size) qs.set('size', String(params.size));
    if (params.user_id) qs.set('user_id', String(params.user_id));
    if (params.asset_id) qs.set('asset_id', String(params.asset_id));
    if (params.keyword) qs.set('keyword', params.keyword);
    return request<AuditPage>(`/audit-logs?${qs.toString()}`);
  },

  // Dashboard
  getDashboardStats: () =>
    request<DashboardStats>('/dashboard/stats'),
};
