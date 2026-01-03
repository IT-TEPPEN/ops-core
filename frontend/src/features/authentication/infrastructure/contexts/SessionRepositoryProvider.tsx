import { useMemo } from "react";
import { SessionLocalStorageRepository } from "../repository";
import { SessionRepositoryContext } from "./SessionRepositoryContext";

export function SessionRepositoryProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const sessionRepository = useMemo(
    () => new SessionLocalStorageRepository(),
    []
  );

  return (
    <SessionRepositoryContext.Provider value={sessionRepository}>
      {children}
    </SessionRepositoryContext.Provider>
  );
}
