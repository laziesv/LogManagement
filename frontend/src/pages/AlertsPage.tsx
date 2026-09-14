import { Bell, Check, ShieldCheck } from "lucide-react";
import { Empty } from "../components/common/Empty";
import { useWorkspaceContext } from "../contexts/WorkspaceContext";
import { displayTime } from "../utils/format";

export function AlertsPage() {
  const {
    user,
    alerts,
    rule,
    setRule,
    actionBusy,
    isAdmin,
    saveRule,
    acknowledge,
  } = useWorkspaceContext();
  return (
    <>
      <section className="card rule-card">
        <div>
          <span className="eyebrow">DETECTION RULE</span>
          <h2>Repeated failed logins</h2>
          <p>
            Group normalized login_failed events by source IP within this
            tenant. Rules apply when new events arrive.
          </p>
        </div>
        {rule && (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              saveRule();
            }}
          >
            <label>
              Failed attempts
              <input
                disabled={!isAdmin}
                type="number"
                min="2"
                max="1000"
                value={rule.threshold}
                onChange={(e) =>
                  setRule({
                    ...rule,
                    threshold: Number(e.target.value),
                  })
                }
              />
            </label>
            <label>
              Window (minutes)
              <input
                disabled={!isAdmin}
                type="number"
                min="1"
                max="60"
                value={rule.window_minutes}
                onChange={(e) =>
                  setRule({
                    ...rule,
                    window_minutes: Number(e.target.value),
                  })
                }
              />
            </label>
            <label className="checkbox">
              <input
                disabled={!isAdmin}
                type="checkbox"
                checked={rule.enabled}
                onChange={(e) =>
                  setRule({ ...rule, enabled: e.target.checked })
                }
              />
              Enabled
            </label>
            {isAdmin && (
              <button className="primary" disabled={actionBusy}>
                Save rule
              </button>
            )}
          </form>
        )}
      </section>
      <section className="card">
        <div className="card-heading">
          <div>
            <h2>Alert history</h2>
            <p>Latest 100 alerts · all time · {user.tenant}</p>
          </div>
        </div>
        {alerts.map((a) => (
          <div className="alert-row" key={a.id}>
            <span className="alert-icon">
              <Bell size={20} />
            </span>
            <div>
              <strong>Repeated failed logins</strong>
              <p>
                <span className="mono">{a.src_ip}</span> · {a.count} attempts ·{" "}
                {displayTime(a.created_at)}
              </p>
            </div>
            <span
              className={`severity ${a.status === "open" ? "medium" : "low"}`}
            >
              {a.status}
            </span>
            {isAdmin && a.status === "open" && (
              <button disabled={actionBusy} onClick={() => acknowledge(a.id)}>
                <Check size={15} />
                Acknowledge
              </button>
            )}
          </div>
        ))}
        {!alerts.length && (
          <Empty
            icon={ShieldCheck}
            title="No alerts yet"
            text="New matching events will appear here when the detection threshold is reached."
          />
        )}
      </section>
    </>
  );
}
