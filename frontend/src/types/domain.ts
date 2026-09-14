export interface User {
  id: string;
  email: string;
  tenant: string;
  role: "admin" | "viewer";
}
export interface LogEvent {
  id: number;
  "@timestamp": string;
  tenant: string;
  source: string;
  event_type: string;
  severity: number;
  src_ip: string;
  user: string;
  host: string;
  action: string;
  fields: Record<string, unknown>;
  raw: unknown;
}
export interface Count {
  name: string;
  count: number;
}
export interface Stats {
  total: number;
  critical: number;
  sources: number;
  alerts: number;
  timeline: Count[];
  top_ips: Count[];
  top_users: Count[];
  top_events: Count[];
}
export interface Alert {
  id: number;
  tenant: string;
  src_ip: string;
  count: number;
  created_at: string;
  status: string;
}
export interface Rule {
  enabled: boolean;
  threshold: number;
  window_minutes: number;
}
