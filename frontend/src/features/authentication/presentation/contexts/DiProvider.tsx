import {
  AuthenticationServiceProvider,
  SessionRepositoryProvider,
} from "../../infrastructure/contexts";
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
          <StartLoginProcessProvider>{children}</StartLoginProcessProvider>
        </AuthenticationServiceProvider>
      </SessionRepositoryProvider>
    </>
  );
}
