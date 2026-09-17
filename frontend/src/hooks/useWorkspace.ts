import { useCallback, useEffect, useRef, useState } from "react";
import { sources } from "../constants/sources";
import { alertsApi } from "../services/alerts";
import { authApi } from "../services/auth";
import { APIError } from "../services/client";
import { logsApi } from "../services/logs";
import type { Alert, LogEvent, Rule, Stats, User } from "../types/domain";
import type { Tab } from "../types/navigation";
import { errorText } from "../utils/errors";
import { abnormalSamples, sample } from "../utils/samples";
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
    controller = useRef<AbortController | null>(null),
    actionInProgress = useRef(false);
  const isAdmin = user.role === "admin";
  const load = useCallback(async (options?: { silent?: boolean }) => {
    if (
      options?.silent &&
      (controller.current || actionInProgress.current || document.hidden)
    )
      return;
    controller.current?.abort();
    const ac = new AbortController();
    controller.current = ac;
    if (!options?.silent) setBusy(true);
    if (!options?.silent) setError("");
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
      controller.current = null;
      if (!options?.silent) {
        setError("Choose a valid time range of up to 31 days.");
        setBusy(false);
      }
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
    const timeout = window.setTimeout(
      () => ac.abort(new DOMException("Request timed out", "TimeoutError")),
      30000,
    );
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
      setError("");
      setUpdated(new Date().toLocaleTimeString());
    } catch (e) {
      if (ac.signal.aborted && ac.signal.reason?.name !== "TimeoutError") return;
      if (ac.signal.reason?.name === "TimeoutError")
        setError("Loading took too long. Check the connection and try Refresh.");
      else if (e instanceof APIError && e.status === 401) onLogout();
      else setError(errorText(e));
    } finally {
      window.clearTimeout(timeout);
      if (controller.current === ac) {
        controller.current = null;
        setBusy(false);
      }
    }
  }, [source, hours, customFrom, customTo, search, offset, onLogout]);
  useEffect(() => {
    void load();
    return () => controller.current?.abort();
  }, [load]);
  useEffect(() => {
    const id = window.setInterval(() => void load({ silent: true }), 10000);
    return () => window.clearInterval(id);
  }, [load]);
  async function action(fn: () => Promise<unknown>, message: string) {
    controller.current?.abort();
    actionInProgress.current = true;
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
      actionInProgress.current = false;
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
      `${records.length} sample events ingested. Auto refresh uses the selected time range.`,
    );
  }

  function demoAbnormal() {
    if (!isAdmin || !rule?.enabled || actionBusy) return;
    const used = new Set(alerts.map((alert) => alert.src_ip));
    const ip = Array.from({ length: 254 }, (_, i) => `192.0.2.${i + 1}`)
      .find((candidate) => !used.has(candidate))!;
    const records = abnormalSamples(rule.threshold, ip);
    void action(
      () => logsApi.ingest(records),
      `${records.length} suspicious login samples sent from ${ip}. Open Alerts to inspect the result.`,
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
    demoAbnormal,
    logout,
    saveRule,
    acknowledge,
  };
}
export type WorkspaceState = ReturnType<typeof useWorkspace>;
