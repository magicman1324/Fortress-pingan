export interface User {
  id: number;
  username: string;
  role: string;
}

export interface Asset {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Session {
  id: number;
  user_id: number;
  asset_id: number;
  status: string;
  client_ip: string;
  started_at: string;
  ended_at: string | null;
  username: string;
  asset_name: string;
}

export interface AuditLog {
  id: number;
  session_id: number;
  user_id: number;
  asset_id: number;
  command: string;
  executed_at: string;
  username: string;
  asset_name: string;
}

export interface AuditPage {
  items: AuditLog[];
  total: number;
  page: number;
  size: number;
}

export interface DashboardStats {
  active_sessions: number;
  total_sessions: number;
  today_audit_count: number;
  recent_audits: AuditLog[];
}

export interface LoginResponse {
  token: string;
  user_id: number;
  username: string;
  role: string;
}
