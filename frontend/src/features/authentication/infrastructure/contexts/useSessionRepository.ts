import { useContext } from "react";
import { SessionRepositoryContext } from "./SessionRepositoryContext";

export function useSessionRepository() {
  const sessionRepository = useContext(SessionRepositoryContext);

  if (!sessionRepository) {
    throw new Error(
      "useSessionRepository must be used within a SessionRepositoryProvider"
    );
  }

  return sessionRepository;
}
