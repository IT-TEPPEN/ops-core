import { createContext, useContext, useMemo, ReactNode } from "react";
import type { ExecutionQueryService } from "../../application";
import { HttpExecutionQueryService } from "../../infrastructure";

/**
 * Context for ExecutionQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const ExecutionQueryServiceContext = createContext<ExecutionQueryService | null>(
  null
);

export function ExecutionQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpExecutionQueryService();
  }, []);

  return (
    <ExecutionQueryServiceContext.Provider value={service}>
      {children}
    </ExecutionQueryServiceContext.Provider>
  );
}

/**
 * Hook to access ExecutionQueryService from context.
 */
export function useExecutionQueryService(): ExecutionQueryService {
  const context = useContext(ExecutionQueryServiceContext);
  if (!context) {
    throw new Error(
      "useExecutionQueryService must be used within ExecutionQueryServiceProvider"
    );
  }
  return context;
}
