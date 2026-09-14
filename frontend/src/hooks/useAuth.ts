import { useCallback, useEffect, useState } from "react";
import { authApi } from "../services/auth";
import { APIError } from "../services/client";
import type { User } from "../types/domain";
import { errorText } from "../utils/errors";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [starting, setStarting] = useState(true);
  const [connectionError, setConnectionError] = useState("");
  const check = useCallback(async () => {
    setStarting(true);
    setConnectionError("");
    try {
      setUser(await authApi.me());
    } catch (e) {
      if (!(e instanceof APIError && e.status === 401))
        setConnectionError(errorText(e));
    } finally {
      setStarting(false);
    }
  }, []);
  useEffect(() => {
    void check();
  }, [check]);

  const onLogout = useCallback(() => setUser(null), []);
  return { user, setUser, starting, connectionError, check, onLogout };
}
