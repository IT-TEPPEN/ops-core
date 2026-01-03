import { useContext } from "react";
import { StartLoginProcessContext } from "./StartLoginProcessContext";

export function useStartLoginProcessUsecase() {
  const usecase = useContext(StartLoginProcessContext);

  if (!usecase) {
    throw new Error(
      "useStartLoginProcessUsecase must be used within a StartLoginProcessProvider"
    );
  }

  return usecase;
}
