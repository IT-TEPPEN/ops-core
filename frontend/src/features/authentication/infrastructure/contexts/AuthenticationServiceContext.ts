import { createContext } from "react";
import { AuthenticationService } from "../../application/services";

export const AuthenticationServiceContext =
  createContext<AuthenticationService | null>(null);
