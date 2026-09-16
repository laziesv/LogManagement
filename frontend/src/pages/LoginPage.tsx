import { ArrowRight, ShieldCheck } from "lucide-react";
import type { FormEvent } from "react";
import { useState } from "react";
import { Brand } from "../components/common/Brand";
import { authApi } from "../services/auth";
import type { User } from "../types/domain";
import { errorText } from "../utils/errors";

export function LoginPage({
  onLogin,
  connectionError,
  retry,
}: {
  onLogin: (u: User) => void;
  connectionError: string;
  retry: () => void;
}) {
  const [email, setEmail] = useState(""),
    [password, setPassword] = useState(""),
    [error, setError] = useState(""),
    [busy, setBusy] = useState(false);
  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError("");
    try {
      onLogin(await authApi.login({ email, password }));
    } catch (e) {
      setError(errorText(e));
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="login-layout">
      <section className="login-story">
        <Brand />
        <div>
          <span className="eyebrow">YOUR ENVIRONMENT, IN FOCUS</span>
          <h1>
            Every event.
            <br />
            One clear picture.
          </h1>
          <p>
            Collect, explore, and understand your logs.
            <br />
            Keep the important signals in sight.
          </p>
          <div className="signal-art" aria-hidden="true">
            {[
              34, 60, 43, 78, 50, 100, 70, 55, 85, 60, 40, 75, 50, 92, 66, 38,
            ].map((h, i) => (
              <span key={i} style={{ height: `${h}%` }} />
            ))}
          </div>
        </div>
        <small>LOGDESK / LOG MANAGEMENT DEMO</small>
      </section>
      <section className="login-form">
        <div>
          <span className="eyebrow">WELCOME BACK</span>
          <h2>Sign in to Logdesk</h2>
          <p className="muted">
            Your workspace is scoped to your organization.
          </p>
          {connectionError && (
            <div className="error" role="alert">
              API unavailable: {connectionError}{" "}
              <button onClick={retry}>Retry</button>
            </div>
          )}
          <form onSubmit={submit}>
            <label>
              Email address
              <input
                type="email"
                autoComplete="username"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </label>
            <label>
              Password
              <input
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>
            {error && (
              <div role="alert" className="error">
                {error}
              </div>
            )}
            <button className="primary" disabled={busy}>
              {busy ? "Signing in…" : "Sign in"}
              <ArrowRight size={16} />
            </button>
          </form>
          <div className="secure-note">
            <ShieldCheck size={18} />
            <span>Admin and Viewer access · Tenant isolation</span>
          </div>
        </div>
      </section>
    </div>
  );
}
