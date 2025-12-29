import { createContext } from "react";
import type { AuthApiAdapter } from "@/shared/api/authApi";

export const AuthApiContext = createContext<AuthApiAdapter | undefined>(
  undefined
);
