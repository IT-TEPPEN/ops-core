import { createContext, ReactNode, useMemo } from "react";
import type { RepositoryCommandService } from "../../application";
import { HttpRepositoryCommandService } from "../../infrastructure";

const RepositoryCommandServiceContext =
  createContext<RepositoryCommandService | null>(null);

export function RepositoryCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpRepositoryCommandService();
  }, []);

  return (
    <RepositoryCommandServiceContext.Provider value={service}>
      {children}
    </RepositoryCommandServiceContext.Provider>
  );
}

export { RepositoryCommandServiceContext };
