import {
  AuthenticationServiceProvider,
  SessionRepositoryProvider,
} from "../../infrastructure/contexts";
import { GetUserIdentityProvider } from "./GetUserIdentityProvider";
import { StartLoginProcessProvider } from "./StartLoginProcessProvider";

export function AuthenticationDiProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      {/* Infrastructure Layer Providers */}
      <SessionRepositoryProvider>
        <AuthenticationServiceProvider>
          {/* Usecase Layer Providers */}
          <GetUserIdentityProvider>
            <StartLoginProcessProvider>{children}</StartLoginProcessProvider>
          </GetUserIdentityProvider>
        </AuthenticationServiceProvider>
      </SessionRepositoryProvider>
    </>
  );
}
