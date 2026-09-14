import type { Alert, Rule } from "../types/domain";
import { api } from "./client";

export const alertsApi = {
  list: (signal?: AbortSignal) => api<Alert[]>("/alerts", { signal }),
  rule: (signal?: AbortSignal) => api<Rule>("/rule", { signal }),
  saveRule: (rule: Rule) =>
    api<Rule>("/rule", { method: "PUT", body: JSON.stringify(rule) }),
  acknowledge: (id: number) =>
    api<void>(`/alerts/${id}/acknowledge`, { method: "POST" }),
};
