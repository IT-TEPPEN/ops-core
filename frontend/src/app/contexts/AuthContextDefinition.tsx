import { createContext } from "react";

interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

export interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: User, refreshToken: string) => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);
