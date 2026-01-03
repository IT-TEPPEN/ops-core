import { createContext } from "react";
import { SessionRepository } from "../../domain/repository";

export const SessionRepositoryContext = createContext<SessionRepository | null>(
  null
);
