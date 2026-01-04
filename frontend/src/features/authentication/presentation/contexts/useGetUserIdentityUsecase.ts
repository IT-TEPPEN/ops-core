import { useContext } from "react";
import { GetUserIdentityContext } from "./GetUserIdentityContext";

export function useGetUserIdentityUsecase() {
  const context = useContext(GetUserIdentityContext);

  if (context === null) {
    throw new Error(
      "useGetUserIdentityUsecase must be used within a GetUserIdentityProvider"
    );
  }

  return context;
}
