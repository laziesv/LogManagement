import { ChevronRight, ShieldCheck } from "lucide-react";
import { nav } from "../../constants/navigation";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";

export function Topbar() {
  const { user, tab } = useWorkspaceContext();
  return (
    <header className="topbar">
      <span>
        Workspace <ChevronRight size={13} />
        <strong>{nav.find((n) => n.id === tab)?.label}</strong>
      </span>
      <div>
        <span className="tenant-badge">
          <ShieldCheck size={13} />
          {user.tenant}
        </span>
        <span className="role-badge">{user.role}</span>
      </div>
    </header>
  );
}
