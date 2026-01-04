import { GetUserIdentityUsecaseImpl } from "../../application/usecase";
import { useAuthenticationService } from "../../infrastructure/contexts";
import { GetUserIdentityContext } from "./GetUserIdentityContext";

export function GetUserIdentityProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const authenticationService = useAuthenticationService();
  const getUserIdentityUsecase = new GetUserIdentityUsecaseImpl(
    authenticationService
  );

  return (
    <GetUserIdentityContext.Provider value={getUserIdentityUsecase}>
      {children}
    </GetUserIdentityContext.Provider>
  );
}
