import { useCallback, useEffect, useRef, useState } from "react";
import { sources } from "../constants/sources";
import { alertsApi } from "../services/alerts";
import { authApi } from "../services/auth";
import { APIError } from "../services/client";
import { logsApi } from "../services/logs";
import type { Alert, LogEvent, Rule, Stats, User } from "../types/domain";
import type { Tab } from "../types/navigation";
import { errorText } from "../utils/errors";
import { sample } from "../utils/samples";
export function useWorkspace(user: User, onLogout: () => void) {
  const [tab, setTab] = useState<Tab>("overview"),
    [source, setSource] = useState(""),
    [hours, setHours] = useState("24"),
    [query, setQuery] = useState(""),
    [search, setSearch] = useState(""),
    [offset, setOffset] = useState(0);
  const [customFrom, setCustomFrom] = useState(""),
    [customTo, setCustomTo] = useState("");
  const [logs, setLogs] = useState<LogEvent[]>([]),
    [total, setTotal] = useState(0),
    [stats, setStats] = useState<Stats | null>(null),
    [alerts, setAlerts] = useState<Alert[]>([]),
    [rule, setRule] = useState<Rule | null>(null);
  const [busy, setBusy] = useState(true),
    [actionBusy, setActionBusy] = useState(false),
    [error, setError] = useState(""),
    [notice, setNotice] = useState(""),
    [selected, setSelected] = useState<LogEvent | null>(null),
    [updated, setUpdated] = useState("");
  const fileInput = useRef<HTMLInputElement>(null),
    controller = useRef<AbortController | null>(null);
  const isAdmin = user.role === "admin";
  const load = useCallback(async () => {
    controller.current?.abort();
    const ac = new AbortController();
    controller.current = ac;
    setBusy(true);
    setError("");
    const to = hours === "custom" ? new Date(customTo) : new Date(),
      from =
        hours === "custom"
          ? new Date(customFrom)
          : new Date(to.getTime() - Number(hours) * 3600000);
    if (
      isNaN(from.getTime()) ||
      isNaN(to.getTime()) ||
      from > to ||
      to.getTime() - from.getTime() > 31 * 86400000
    ) {
      setError("Choose a valid time range of up to 31 days.");
      setBusy(false);
      return;
    }
    const params = new URLSearchParams({
      source,
      q: search,
      from: from.toISOString(),
      to: to.toISOString(),
      limit: "50",
      offset: String(offset),
    });
    try {
      const [l, s, a, r] = await Promise.all([
        logsApi.list(params, ac.signal),
        logsApi.stats(params, ac.signal),
        alertsApi.list(ac.signal),
        alertsApi.rule(ac.signal),
      ]);
      if (ac.signal.aborted) return;
      setLogs(l.items);
      setTotal(l.total);
      setStats(s);
      setAlerts(a);
      setRule(r);
      setUpdated(new Date().toLocaleTimeString());
    } catch (e) {
      if (ac.signal.aborted) return;
      if (e instanceof APIError && e.status === 401) onLogout();
      else setError(errorText(e));
    } finally {
      if (!ac.signal.aborted) setBusy(false);
    }
  }, [source, hours, customFrom, customTo, search, offset, onLogout]);
  useEffect(() => {
    void load();
    return () => controller.current?.abort();
  }, [load]);
  async function action(fn: () => Promise<unknown>, message: string) {
    setActionBusy(true);
    setError("");
    setNotice("");
    try {
      await fn();
      setNotice(message);
      await load();
    } catch (e) {
      setError(errorText(e));
    } finally {
      setActionBusy(false);
    }
  }
  async function upload(file?: File) {
    if (!file) return;
    await action(
      () => logsApi.upload(file),
      `${file.name} imported successfully.`,
    );
    if (fileInput.current) fileInput.current.value = "";
  }
  function demo(source?: string) {
    const records = source
      ? [sample(source)]
      : [
          ...sources.map(sample),
          ...Array.from({ length: 5 }, () => ({
            ...sample("api"),
            event_type: "login_failed",
            severity: 5,
          })),
        ];
    void action(
      () => logsApi.ingest(records),
      `${records.length} sample events ingested. Refresh uses the selected time range.`,
    );
  }

  async function logout() {
    setActionBusy(true);
    setError("");
    try {
      await authApi.logout();
      onLogout();
    } catch (e) {
      setError(errorText(e));
    } finally {
      setActionBusy(false);
    }
  }
  function saveRule() {
    if (rule)
      void action(() => alertsApi.saveRule(rule), "Detection rule saved.");
  }
  function acknowledge(id: number) {
    void action(() => alertsApi.acknowledge(id), "Alert acknowledged.");
  }

  return {
    user,
    onLogout,
    tab,
    setTab,
    source,
    setSource,
    hours,
    setHours,
    query,
    setQuery,
    search,
    setSearch,
    offset,
    setOffset,
    customFrom,
    setCustomFrom,
    customTo,
    setCustomTo,
    logs,
    total,
    stats,
    alerts,
    rule,
    setRule,
    busy,
    actionBusy,
    error,
    notice,
    setNotice,
    selected,
    setSelected,
    updated,
    fileInput,
    isAdmin,
    load,
    action,
    upload,
    demo,
    logout,
    saveRule,
    acknowledge,
  };
}
export type WorkspaceState = ReturnType<typeof useWorkspace>;
