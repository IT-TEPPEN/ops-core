import { useContext } from "react";
import { AuthApiContext } from "@/app/contexts/AuthApiContext";
import type { AuthApiAdapter } from "@/shared/api/authApi";

export function useAuthApi(): AuthApiAdapter {
  const context = useContext(AuthApiContext);
  if (context === undefined) {
    throw new Error("useAuthApi must be used within an AuthApiProvider");
  }
  return context;
}
