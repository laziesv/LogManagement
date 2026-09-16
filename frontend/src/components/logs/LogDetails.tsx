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
  const fields = event.fields ?? {};
  const cloud = fields.cloud && typeof fields.cloud === "object" && !Array.isArray(fields.cloud)
    ? fields.cloud as Record<string, unknown>
    : {};
  const extraFields: [string, unknown][] = [
    ["Host", event.host], ["Action", event.action], ["Outcome", fields.outcome],
    ["Vendor", fields.vendor], ["Product", fields.product],
    ["Source port", fields.src_port], ["Destination IP", fields.dst_ip],
    ["Destination port", fields.dst_port], ["Protocol", fields.protocol],
    ["URL / path", fields.url], ["HTTP method", fields.http_method],
    ["HTTP status", fields.status_code], ["Cloud account", cloud.account_id],
    ["Cloud region", cloud.region], ["Cloud service", cloud.service],
  ];
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
        <dl>
          {extraFields.filter(([, value]) =>
            (typeof value === "string" && value !== "") || typeof value === "number"
          ).map(([label, value]) => (
            <div key={label}>
              <dt>{label}</dt>
              <dd>{String(value)}</dd>
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
