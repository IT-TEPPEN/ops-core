import { createContext, useContext, useMemo, ReactNode } from "react";
import type { GroupQueryService } from "../../application";
import { HttpGroupQueryService } from "../../infrastructure";

/**
 * Context for GroupQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const GroupQueryServiceContext = createContext<GroupQueryService | null>(null);

export function GroupQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpGroupQueryService();
  }, []);

  return (
    <GroupQueryServiceContext.Provider value={service}>
      {children}
    </GroupQueryServiceContext.Provider>
  );
}

/**
 * Hook to access GroupQueryService from context.
 */
export function useGroupQueryService(): GroupQueryService {
  const context = useContext(GroupQueryServiceContext);
  if (!context) {
    throw new Error(
      "useGroupQueryService must be used within GroupQueryServiceProvider"
    );
  }
  return context;
}
