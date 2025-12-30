import { createContext } from "react";

interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

interface Identity {
  id: string;
  provider: string;
  email: string;
  name: string;
  linkedAt: string;
  lastUsedAt: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  identities: Identity[];
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: User, refreshToken: string) => void;
  logout: () => void;
  linkProvider: (provider: string) => void;
  unlinkIdentity: (identityId: string) => Promise<void>;
  refreshIdentities: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);
