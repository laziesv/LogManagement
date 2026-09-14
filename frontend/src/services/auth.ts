import type { User } from "../types/domain";
import { api } from "./client";

export const authApi = {
  me: () => api<User>("/auth/me"),
  login: (credentials: { email: string; password: string }) =>
    api<User>("/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    }),
  logout: () => api<void>("/auth/logout", { method: "POST" }),
};
