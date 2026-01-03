import { useMemo } from "react";
import { AuthenticationServiceContext } from "./AuthenticationServiceContext";
import { SsoService } from "../services";

export function AuthenticationServiceProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const services = useMemo(() => new SsoService(), []);

  return (
    <AuthenticationServiceContext.Provider value={services}>
      {children}
    </AuthenticationServiceContext.Provider>
  );
}
