import { Check, X } from "lucide-react";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";

export function Notifications() {
  const { error, notice, setNotice, updated } = useWorkspaceContext();
  return (
    <>
      {error && (
        <div className="error" role="alert">
          {error}
          {updated && " Previously loaded data may be stale."}
        </div>
      )}
      {notice && (
        <div className="notice" role="status">
          <Check size={16} />
          {notice}
          <button
            aria-label="Dismiss notification"
            onClick={() => setNotice("")}
          >
            <X size={15} />
          </button>
        </div>
      )}
    </>
  );
}
