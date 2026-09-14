import {
  Activity,
  LayoutDashboard,
  FileText,
  Bell,
  Database,
} from "lucide-react";
import type { Tab } from "../types/navigation";

export const titles: Record<Tab, string> = {
  overview: "Workspace overview",
  logs: "Log explorer",
  alerts: "Alerts & detection",
  sources: "Data sources",
};
export const descriptions: Record<Tab, string> = {
  overview: "A clearer view of everything happening in your environment.",
  logs: "Search, inspect, and trace events across your connected sources.",
  alerts: "Turn repeated signals into something you can act on.",
  sources: "Bring your events together in one consistent format.",
};

export const nav: { id: Tab; label: string; icon: typeof Activity }[] = [
  { id: "overview", label: "Overview", icon: LayoutDashboard },
  { id: "logs", label: "Log explorer", icon: FileText },
  { id: "alerts", label: "Alerts", icon: Bell },
  { id: "sources", label: "Data sources", icon: Database },
];
