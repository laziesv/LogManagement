import { Layers3 } from "lucide-react";
import { useAuth } from "./hooks/useAuth";
import { Workspace } from "./layouts/Workspace";
import { LoginPage } from "./pages/LoginPage";

export default function App() {
  const { user, setUser, starting, connectionError, check, onLogout } =
    useAuth();
  if (starting)
    return (
      <div className="splash">
        <Layers3 size={38} />
        <p>Opening your workspace…</p>
      </div>
    );
  if (!user)
    return (
      <LoginPage
        onLogin={setUser}
        connectionError={connectionError}
        retry={check}
      />
    );
  return <Workspace user={user} onLogout={onLogout} />;
}
