import { LogOut, Radio, ShieldCheck } from "lucide-react";
import { nav } from "../../constants/navigation";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";
import { Brand } from "../common/Brand";

export function Sidebar() {
  const { user, tab, setTab, stats, isAdmin, logout } = useWorkspaceContext();
  return (
    <aside className="sidebar">
      <Brand />
      <div className="workspace-pill">
        <span className="tenant-avatar">D</span>
        <div>
          Demo workspace<small>{user.tenant}</small>
        </div>
        <ShieldCheck size={16} />
      </div>
      <span className="nav-label">WORKSPACE</span>
      <nav>
        {nav.map((n) => (
          <button
            key={n.id}
            className={tab === n.id ? "active" : ""}
            onClick={() => setTab(n.id)}
          >
            <n.icon size={18} />
            {n.label}
            {n.id === "alerts" && !!stats?.alerts && (
              <span className="nav-count">{stats.alerts}</span>
            )}
          </button>
        ))}
      </nav>
      <div className="sidebar-note">
        <Radio size={20} />
        <strong>A place for every signal.</strong>
        <p>Explore events from your network, cloud, and applications.</p>
      </div>
      <div className="profile">
        <span className="profile-avatar">{isAdmin ? "A" : "V"}</span>
        <div>
          <strong>{isAdmin ? "Administrator" : "Viewer"}</strong>
          <small title={user.email}>{user.email}</small>
        </div>
        <button
          className="icon-button"
          title="Sign out"
          aria-label="Sign out"
          onClick={() => void logout()}
        >
          <LogOut size={17} />
        </button>
      </div>
    </aside>
  );
}
