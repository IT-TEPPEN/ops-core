import { useContext } from "react";
import { SignOutContext } from "./SignOutContext";

export function useSignOutUsecase() {
  const context = useContext(SignOutContext);

  if (context === null) {
    throw new Error(
      "useSignOutUsecase must be used within a SignOutProvider"
    );
  }

  return context;
}
