import { X } from "lucide-react";
import { useEffect, useRef } from "react";
import type { LogEvent } from "../../types/domain";

export function LogDetails({
  event,
  close,
}: {
  event: LogEvent;
  close: () => void;
}) {
  const ref = useRef<HTMLDialogElement>(null);
  useEffect(() => {
    const d = ref.current;
    d?.showModal();
    return () => d?.close();
  }, []);
  return (
    <dialog
      ref={ref}
      className="log-dialog"
      onCancel={close}
      onClick={(e) => {
        if (e.target === e.currentTarget) close();
      }}
    >
      <div className="card-heading">
        <div>
          <span className="eyebrow">EVENT #{event.id}</span>
          <h2>{event.event_type}</h2>
        </div>
        <button aria-label="Close details" onClick={close}>
          <X size={18} />
        </button>
      </div>
      <div className="detail-body">
        <dl>
          {[
            ["Timestamp", event["@timestamp"]],
            ["Tenant", event.tenant],
            ["Source", event.source],
            ["Source IP", event.src_ip || "—"],
            ["Severity", event.severity],
            ["User", event.user || "—"],
          ].map(([k, v]) => (
            <div key={k}>
              <dt>{k}</dt>
              <dd>{v}</dd>
            </div>
          ))}
        </dl>
        <h3>Normalized fields</h3>
        <pre>{JSON.stringify(event.fields, null, 2)}</pre>
        <h3>Original payload</h3>
        <pre>
          {typeof event.raw === "string"
            ? event.raw
            : JSON.stringify(event.raw, null, 2)}
        </pre>
      </div>
    </dialog>
  );
}
