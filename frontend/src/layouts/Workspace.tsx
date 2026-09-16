import { Notifications } from "../components/layout/Notifications";
import { PageHeader } from "../components/layout/PageHeader";
import { Sidebar } from "../components/layout/Sidebar";
import { StatusFooter } from "../components/layout/StatusFooter";
import { Topbar } from "../components/layout/Topbar";
import { LogDetails } from "../components/logs/LogDetails";
import { LogFilters } from "../components/logs/LogFilters";
import { useWorkspace } from "../hooks/useWorkspace";
import { AlertsPage } from "../pages/AlertsPage";
import { LogsPage } from "../pages/LogsPage";
import { OverviewPage } from "../pages/OverviewPage";
import { SourcesPage } from "../pages/SourcesPage";
import type { User } from "../types/domain";

import { WorkspaceProvider } from "../contexts/WorkspaceContext";
export function Workspace({
  user,
  onLogout,
}: {
  user: User;
  onLogout: () => void;
}) {
  const state = useWorkspace(user, onLogout);
  return (
    <WorkspaceProvider value={state}>
      <div className="app-shell">
        <Sidebar />
        <main>
          <Topbar />
          <div className="page" key={state.tab}>
            <PageHeader />
            <Notifications />
            {(state.tab === "overview" || state.tab === "logs") && (
              <LogFilters />
            )}
            {state.tab === "overview" && <OverviewPage />}
            {state.tab === "logs" && <LogsPage />}
            {state.tab === "alerts" && <AlertsPage />}
            {state.tab === "sources" && <SourcesPage />}
            <StatusFooter />
          </div>
        </main>
        {state.selected && (
          <LogDetails
            event={state.selected}
            close={() => state.setSelected(null)}
          />
        )}
      </div>
    </WorkspaceProvider>
  );
}
