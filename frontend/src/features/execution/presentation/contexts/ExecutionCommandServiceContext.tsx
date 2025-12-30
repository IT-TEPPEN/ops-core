import { createContext, useContext, useMemo, ReactNode } from "react";
import type { ExecutionCommandService } from "../../application";
import { HttpExecutionCommandService } from "../../infrastructure";

/**
 * Context for ExecutionCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const ExecutionCommandServiceContext = createContext<ExecutionCommandService | null>(
  null
);

export function ExecutionCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpExecutionCommandService();
  }, []);

  return (
    <ExecutionCommandServiceContext.Provider value={service}>
      {children}
    </ExecutionCommandServiceContext.Provider>
  );
}

/**
 * Hook to access ExecutionCommandService from context.
 */
export function useExecutionCommandService(): ExecutionCommandService {
  const context = useContext(ExecutionCommandServiceContext);
  if (!context) {
    throw new Error(
      "useExecutionCommandService must be used within ExecutionCommandServiceProvider"
    );
  }
  return context;
}
