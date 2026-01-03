import { createContext } from "react";
import { StartLoginProcessUsecase } from "../../application/usecase/StartLoginProcess";

export const StartLoginProcessContext =
  createContext<StartLoginProcessUsecase | null>(null);
