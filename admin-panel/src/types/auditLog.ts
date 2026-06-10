// Audit Log типы
export interface AuditLog {
  id: string;
  admin_id: string | null;
  admin_username: string | null;
  action: string;
  table_name: string | null;
  record_id: string | null;
  details: Record<string, unknown> | null;
  created_at: string;
}
