import { Bell, Layers3, Radio, ShieldCheck } from "lucide-react";
import { Ranking } from "../components/dashboard/Ranking";
import { Timeline } from "../components/dashboard/Timeline";
import { LogTable } from "../components/logs/LogTable";
import { useWorkspaceContext } from "../contexts/WorkspaceContext";
import { number } from "../utils/format";

export function OverviewPage() {
  const { stats, busy } = useWorkspaceContext();
  return (
    <>
      <>
        <div className="metric-grid">
          {[
            {
              label: "Total events",
              value: stats?.total,
              icon: Layers3,
              note: "In selected range",
            },
            {
              label: "High severity",
              value: stats?.critical,
              icon: ShieldCheck,
              note: "Severity 8–10",
            },
            {
              label: "Active sources",
              value: stats?.sources,
              icon: Radio,
              note: "In selected range",
            },
            {
              label: "Open alerts",
              value: stats?.alerts,
              icon: Bell,
              note: "All time · this tenant",
            },
          ].map((m) => (
            <div className="metric card" key={m.label}>
              <div>
                {m.label}
                <m.icon size={18} />
              </div>
              <strong>{m.value === undefined ? "—" : number(m.value)}</strong>
              <small>{m.note}</small>
            </div>
          ))}
        </div>
        <div className="card chart-card">
          <div className="card-heading">
            <div>
              <h2>Event activity _TH</h2>
              <p>Hourly event volume · UTC</p>
            </div>
            <span className="chart-legend">
              <i />
              Events
            </span>
          </div>
          <Timeline points={stats?.timeline || []} loading={busy && !stats} />
        </div>
        <div className="top-grid">
          <Ranking title="Top source IPs" data={stats?.top_ips || []} />
          <Ranking title="Top users" data={stats?.top_users || []} />
          <Ranking title="Top event types" data={stats?.top_events || []} />
        </div>
      </>
      <LogTable />
    </>
  );
}
