import {
  AuthenticationServiceProvider,
  SessionRepositoryProvider,
} from "../../infrastructure/contexts";
import { GetUserIdentityProvider } from "./GetUserIdentityProvider";
import { SignOutProvider } from "./SignOutProvider";
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
            <SignOutProvider>
              <StartLoginProcessProvider>{children}</StartLoginProcessProvider>
            </SignOutProvider>
          </GetUserIdentityProvider>
        </AuthenticationServiceProvider>
      </SessionRepositoryProvider>
    </>
  );
}
