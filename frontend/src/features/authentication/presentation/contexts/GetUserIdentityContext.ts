import { createContext } from "react";
import { GetUserIdentityUsecase } from "../../application/usecase";

export const GetUserIdentityContext =
  createContext<GetUserIdentityUsecase | null>(null);
