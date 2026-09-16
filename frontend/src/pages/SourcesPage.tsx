import {
  ArrowDownToLine,
  ArrowRight,
  Check,
  Database,
  Plus,
  Upload,
  ShieldAlert,
} from "lucide-react";
import { sources } from "../constants/sources";
import { useWorkspaceContext } from "../contexts/WorkspaceContext";

export function SourcesPage() {
  const { actionBusy, fileInput, isAdmin, demo, demoAbnormal, rule } = useWorkspaceContext();
  return (
    <>
      <div className="source-banner">
        <div>
          <span className="eyebrow">BRING YOUR DATA IN</span>
          <h2>Seven source types. One shared schema.</h2>
          <p>
            Use HTTP JSON, file imports, or the Syslog collector. Sample events
            are written to your tenant's database.
          </p>
        </div>
        {isAdmin && (
          <button
            className="primary"
            disabled={actionBusy}
            onClick={() => demo()}
          >
            <Plus size={16} />
            Send demo batch
          </button>
        )}
      </div>
      {isAdmin && (
        <section className="card import-card" style={{ marginTop: 0, marginBottom: 22 }}>
          <ShieldAlert size={28} />
          <div>
            <h2>Simulate suspicious logins</h2>
            <p>
              Send {rule?.threshold ?? 5} failed logins from one sample IP to test your alert rule.
              {rule?.enabled
                ? " Then open Alerts to inspect the notification."
                : " Enable the failed-login rule in Alerts first."}
            </p>
          </div>
          <button onClick={demoAbnormal} disabled={actionBusy || !rule?.enabled}>
            <ShieldAlert size={16} />
            Send abnormal sample
          </button>
        </section>
      )}
      <section className="card import-card" style={{ marginTop: 0, marginBottom: 22 }}>
        <Database size={28} />
        <div>
          <h2>Real web access logs</h2>
          <p>
            Nginx keeps web access logs local during the demo, so refreshing the
            website will not create new application logs.
          </p>
          <p>Use the Syslog UDP/TCP collector on port 5514 when you want to test live ingestion.</p>
        </div>
      </section>
      <div className="source-grid">
        {sources.map((s, i) => (
          <article className="card source-card" key={s}>
            <div className={`source-icon source-bg-${i % 3}`}>
              <Database size={22} />
            </div>
            <h2>
              {
                (
                  {
                    api: "HTTP API",
                    firewall: "Firewall",
                    network: "Network",
                    crowdstrike: "CrowdStrike",
                    aws: "AWS CloudTrail",
                    m365: "Microsoft 365",
                    ad: "Windows / AD",
                  } as Record<string, string>
                )[s]
              }
            </h2>
            <p>
              {s === "firewall" || s === "network"
                ? "Syslog UDP/TCP or JSON ingest"
                : "JSON API or batch file import"}
            </p>
            <span className="supported">
              <Check size={13} />
              Parser available
            </span>
            {isAdmin && (
              <button disabled={actionBusy} onClick={() => demo(s)}>
                Send sample
                <ArrowRight size={15} />
              </button>
            )}
          </article>
        ))}
      </div>
      <section className="card import-card">
        <Upload size={28} />
        <div>
          <h2>Have a log file ready?</h2>
          <p>
            Import a JSON object, array, or AWS Records envelope. Up to 1,000
            events / 2 MB per request.
          </p>
        </div>
        {isAdmin && (
          <button
            onClick={() => fileInput.current?.click()}
            disabled={actionBusy}
          >
            <ArrowDownToLine size={16} />
            Choose JSON file
          </button>
        )}
      </section>
    </>
  );
}
