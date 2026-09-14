import { Plus, RefreshCw } from "lucide-react";
import { descriptions, titles } from "../../constants/navigation";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";

export function PageHeader() {
  const { tab, busy, actionBusy, fileInput, isAdmin, load, upload } =
    useWorkspaceContext();
  return (
    <div className="page-heading">
      <div>
        <span className="eyebrow">OBSERVE. UNDERSTAND. ACT.</span>
        <h1>{titles[tab]}</h1>
        <p>{descriptions[tab]}</p>
      </div>
      <div className="heading-actions">
        <button disabled={busy || actionBusy} onClick={() => void load()}>
          <RefreshCw size={15} className={busy ? "spin" : ""} />
          Refresh
        </button>
        {isAdmin && (
          <button
            className="primary"
            disabled={actionBusy}
            onClick={() => fileInput.current?.click()}
          >
            <Plus size={16} />
            Import logs
          </button>
        )}
        <input
          ref={fileInput}
          type="file"
          accept=".json,application/json"
          className="sr-only"
          aria-label="Import JSON logs"
          onChange={(e) => void upload(e.target.files?.[0])}
        />
      </div>
    </div>
  );
}
