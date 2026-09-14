import type { LogEvent, Stats } from "../types/domain";
import { api } from "./client";

export const logsApi = {
  list: (params: URLSearchParams, signal?: AbortSignal) =>
    api<{ items: LogEvent[]; total: number }>(`/logs?${params}`, { signal }),
  stats: (params: URLSearchParams, signal?: AbortSignal) =>
    api<Stats>(`/stats?${params}`, { signal }),
  ingest: (records: unknown) =>
    api<{ accepted: number; tenant: string }>("/ingest", {
      method: "POST",
      body: JSON.stringify(records),
    }),
  upload: (file: File) => {
    const body = new FormData();
    body.append("file", file);
    return api<{ accepted: number; tenant: string }>("/ingest", {
      method: "POST",
      body,
    });
  },
};
