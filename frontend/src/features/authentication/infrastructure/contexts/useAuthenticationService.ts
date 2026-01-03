import { useContext } from "react";
import { AuthenticationServiceContext } from "./AuthenticationServiceContext";
import { AuthenticationService } from "../../application/services";

export function useAuthenticationService(): AuthenticationService {
  const context = useContext(AuthenticationServiceContext);

  if (!context) {
    throw new Error(
      "useAuthenticationService must be used within an AuthenticationServiceProvider"
    );
  }

  return context;
}
