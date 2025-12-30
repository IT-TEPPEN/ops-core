import { createContext, useContext, useMemo, ReactNode } from "react";
import type { GroupCommandService } from "../../application";
import { HttpGroupCommandService } from "../../infrastructure";

/**
 * Context for GroupCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const GroupCommandServiceContext = createContext<GroupCommandService | null>(
  null
);

export function GroupCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpGroupCommandService();
  }, []);

  return (
    <GroupCommandServiceContext.Provider value={service}>
      {children}
    </GroupCommandServiceContext.Provider>
  );
}

/**
 * Hook to access GroupCommandService from context.
 */
export function useGroupCommandService(): GroupCommandService {
  const context = useContext(GroupCommandServiceContext);
  if (!context) {
    throw new Error(
      "useGroupCommandService must be used within GroupCommandServiceProvider"
    );
  }
  return context;
}
