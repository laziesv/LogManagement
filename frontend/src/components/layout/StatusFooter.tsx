import { useWorkspaceContext } from "../../contexts/WorkspaceContext";

export function StatusFooter() {
  const { busy, error, updated } = useWorkspaceContext();
  return (
    <footer>
      <span>
        <span className={`status-dot ${error ? "warning" : ""}`} />
        {busy
          ? "Updating workspace…"
          : error
            ? "Check connection or filters"
            : updated
              ? `Last updated ${updated}`
              : "Waiting for data"}
      </span>
      <span>Logdesk / Demo environment</span>
    </footer>
  );
}
