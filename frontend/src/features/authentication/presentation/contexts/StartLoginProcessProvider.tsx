import { useMemo } from "react";
import { StartLoginProcessUsecaseImpl } from "../../application/usecase/StartLoginProcess";
import {
  useAuthenticationService,
  useSessionRepository,
} from "../../infrastructure/contexts";
import { StartLoginProcessContext } from "./StartLoginProcessContext";

export function StartLoginProcessProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const authenticationService = useAuthenticationService();
  const sessionRepository = useSessionRepository();

  const startLoginProcessUsecase = useMemo(
    () =>
      new StartLoginProcessUsecaseImpl(
        authenticationService,
        sessionRepository
      ),
    [authenticationService, sessionRepository]
  );

  return (
    <StartLoginProcessContext.Provider value={startLoginProcessUsecase}>
      {children}
    </StartLoginProcessContext.Provider>
  );
}
