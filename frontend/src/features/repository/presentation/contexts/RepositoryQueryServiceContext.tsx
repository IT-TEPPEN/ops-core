import { createContext, ReactNode, useMemo } from "react";
import type { RepositoryQueryService } from "../../application";
import { HttpRepositoryQueryService } from "../../infrastructure";

const RepositoryQueryServiceContext =
  createContext<RepositoryQueryService | null>(null);

export function RepositoryQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpRepositoryQueryService();
  }, []);

  return (
    <RepositoryQueryServiceContext.Provider value={service}>
      {children}
    </RepositoryQueryServiceContext.Provider>
  );
}

export { RepositoryQueryServiceContext };
