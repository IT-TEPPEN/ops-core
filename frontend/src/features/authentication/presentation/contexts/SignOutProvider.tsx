import { SignOutUsecaseImpl } from "../../application/usecase";
import {
  useAuthenticationService,
  useSessionRepository,
} from "../../infrastructure/contexts";
import { SignOutContext } from "./SignOutContext";

export function SignOutProvider({ children }: { children: React.ReactNode }) {
  const sessionRepository = useSessionRepository();
  const authenticationService = useAuthenticationService();
  const signOutUsecase = new SignOutUsecaseImpl(
    sessionRepository,
    authenticationService
  );

  return (
    <SignOutContext.Provider value={signOutUsecase}>
      {children}
    </SignOutContext.Provider>
  );
}
