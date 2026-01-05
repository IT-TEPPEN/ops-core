import { createContext } from "react";
import { SignOutUsecase } from "../../application/usecase";

export const SignOutContext = createContext<SignOutUsecase | null>(null);
