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
    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 60000);
    return api<{ accepted: number; tenant: string }>("/ingest", {
      method: "POST",
      body,
      signal: controller.signal,
    }).catch((error: unknown) => {
      if (controller.signal.aborted)
        throw new Error("Upload timed out. Check Log Explorer before trying again.");
      throw error;
    }).finally(() => window.clearTimeout(timeout));
  },
};
