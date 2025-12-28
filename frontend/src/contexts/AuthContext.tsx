import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";

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

interface AuthContextType {
  user: User | null;
  token: string | null;
  identities: Identity[];
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  linkProvider: (provider: string) => void;
  unlinkIdentity: (identityId: string) => Promise<void>;
  refreshIdentities: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const TOKEN_KEY = "auth_token";
const USER_KEY = "auth_user";

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [identities, setIdentities] = useState<Identity[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize auth state from localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem(TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        setToken(storedToken);
        setUser(parsedUser);
      } catch (error) {
        console.error("Failed to parse stored user data:", error);
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }

    setIsLoading(false);
  }, []);

  // Fetch identities when authenticated
  useEffect(() => {
    if (token) {
      refreshIdentities();
    }
  }, [token]);

  const refreshIdentities = async () => {
    if (!token) return;

    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/auth/identities",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setIdentities(data.identities || []);
      }
    } catch (error) {
      console.error("Failed to fetch identities:", error);
    }
  };

  const login = (newToken: string, newUser: User) => {
    localStorage.setItem(TOKEN_KEY, newToken);
    localStorage.setItem(USER_KEY, JSON.stringify(newUser));
    setToken(newToken);
    setUser(newUser);
  };

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    // Also remove OAuth token if it exists
    localStorage.removeItem("oauth_token");
    setToken(null);
    setUser(null);
    setIdentities([]);
  };

  const linkProvider = async (provider: string) => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/auth/${provider}/login`
      );
      const data = await response.json();

      sessionStorage.setItem(`${provider}_auth_state`, data.state);
      sessionStorage.setItem("link_mode", "true");

      window.location.href = data.auth_url;
    } catch (error) {
      console.error(`Failed to initiate ${provider} login:`, error);
      throw error;
    }
  };

  const unlinkIdentity = async (identityId: string) => {
    if (!token) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/auth/identities/${identityId}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to unlink identity");
      }

      await refreshIdentities();
    } catch (error) {
      console.error("Failed to unlink identity:", error);
      throw error;
    }
  };

  const isAuthenticated = !!token && !!user;

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        identities,
        isAuthenticated,
        isLoading,
        login,
        logout,
        linkProvider,
        unlinkIdentity,
        refreshIdentities,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
