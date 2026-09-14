import type { ReactNode } from "react";
import { createContext, useContext } from "react";
import type { WorkspaceState } from "../hooks/useWorkspace";
const WorkspaceContext = createContext<WorkspaceState | null>(null);
export function WorkspaceProvider({
  value,
  children,
}: {
  value: WorkspaceState;
  children: ReactNode;
}) {
  return (
    <WorkspaceContext.Provider value={value}>
      {children}
    </WorkspaceContext.Provider>
  );
}
export function useWorkspaceContext() {
  const value = useContext(WorkspaceContext);
  if (!value) throw new Error("WorkspaceProvider is required");
  return value;
}
