import { createContext, useContext, useMemo, ReactNode } from "react";
import type { UserQueryService } from "../../application/services";
import { HttpUserQueryService } from "../../infrastructure/services";

/**
 * Context for UserQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const UserQueryServiceContext = createContext<UserQueryService | null>(null);

export function UserQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpUserQueryService();
  }, []);

  return (
    <UserQueryServiceContext.Provider value={service}>
      {children}
    </UserQueryServiceContext.Provider>
  );
}

/**
 * Hook to access UserQueryService from context.
 */
export function useUserQueryService(): UserQueryService {
  const context = useContext(UserQueryServiceContext);
  if (!context) {
    throw new Error(
      "useUserQueryService must be used within UserQueryServiceProvider"
    );
  }
  return context;
}
