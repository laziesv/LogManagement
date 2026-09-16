import { ArrowRight, ChevronLeft, ChevronRight, FileText } from "lucide-react";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";
import { displayTime, number } from "../../utils/format";
import { Empty } from "../common/Empty";

export function LogTable() {
  const { tab, setTab, offset, setOffset, logs, total, busy, updated, setSelected } =
    useWorkspaceContext();
  return (
    <section className="card log-card">
      <div className="card-heading">
        <div>
          <h2>
            {tab === "overview" ? "Recent events" : "Events"}{" "}
            <span className="pill">{number(total)}</span>
          </h2>
          <p>Normalized logs across your selected sources</p>
        </div>
        {tab === "logs" && (
          <div className="live-refresh">
            <span className="status-dot" />
            Live refresh
            {updated && <small>Last updated {updated}</small>}
          </div>
        )}
        {tab === "overview" && (
          <button className="text-button" onClick={() => setTab("logs")}>
            Explore all logs
            <ArrowRight size={15} />
          </button>
        )}
      </div>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>Source</th>
              <th>Event</th>
              <th>Source IP</th>
              <th>Severity</th>
              <th>User</th>
              <th>
                <span className="sr-only">Details</span>
              </th>
            </tr>
          </thead>
          <tbody>
            {(tab === "overview" ? logs.slice(0, 7) : logs).map((l) => (
              <tr key={l.id}>
                <td className="mono time-cell">
                  {displayTime(l["@timestamp"])}
                </td>
                <td>
                  <span className={`source-dot source-${l.source}`} />
                  {l.source}
                </td>
                <td className="event-cell">{l.event_type}</td>
                <td className="mono">{l.src_ip || "—"}</td>
                <td>
                  <span
                    className={`severity ${l.severity >= 8 ? "high" : l.severity >= 5 ? "medium" : "low"}`}
                  >
                    {l.severity >= 8
                      ? "High"
                      : l.severity >= 5
                        ? "Medium"
                        : "Low"}{" "}
                    · {l.severity}
                  </span>
                </td>
                <td>{l.user || "—"}</td>
                <td>
                  <button
                    className="icon-button"
                    aria-label={`Inspect log ${l.id}`}
                    onClick={() => setSelected(l)}
                  >
                    <ArrowRight size={16} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!logs.length && (
        <Empty
          icon={FileText}
          title={busy ? "Loading events…" : "No events in this view"}
          text="Import logs or adjust your source, search, and time filters."
        />
      )}
      {tab === "logs" && (
        <div className="pagination">
          <span>
            {total
              ? `${offset + 1}–${Math.min(offset + 50, total)} of ${number(total)}`
              : "0 events"}
          </span>
          <div>
            <button
              disabled={offset === 0 || busy}
              onClick={() => setOffset(Math.max(0, offset - 50))}
            >
              <ChevronLeft size={15} />
              Previous
            </button>
            <button
              disabled={offset + 50 >= total || busy}
              onClick={() => setOffset(offset + 50)}
            >
              Next
              <ChevronRight size={15} />
            </button>
          </div>
        </div>
      )}
    </section>
  );
}
