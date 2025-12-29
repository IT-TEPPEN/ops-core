import { useReducer, useEffect, ReactNode, useMemo } from "react";
import { AuthContext } from "./AuthContextDefinition";
import { AuthApi } from "@/shared/api/authApi";

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

const TOKEN_KEY = "auth_token";
const USER_KEY = "auth_user";

interface AuthProviderProps {
  children: ReactNode;
}

// State type
interface AuthState {
  user: User | null;
  token: string | null;
  identities: Identity[];
  isLoading: boolean;
}

// Action types
type AuthAction =
  | { type: "SET_TOKEN"; payload: string | null }
  | { type: "SET_USER"; payload: User | null }
  | { type: "SET_IDENTITIES"; payload: Identity[] }
  | { type: "SET_LOADING"; payload: boolean }
  | { type: "LOGIN"; payload: { token: string; user: User } }
  | { type: "LOGOUT" };

// Reducer
function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case "SET_TOKEN":
      return { ...state, token: action.payload };
    case "SET_USER":
      return { ...state, user: action.payload };
    case "SET_IDENTITIES":
      return { ...state, identities: action.payload };
    case "SET_LOADING":
      return { ...state, isLoading: action.payload };
    case "LOGIN":
      return {
        ...state,
        token: action.payload.token,
        user: action.payload.user,
      };
    case "LOGOUT":
      return {
        ...state,
        token: null,
        user: null,
        identities: [],
      };
    default:
      return state;
  }
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [state, dispatch] = useReducer(authReducer, {
    user: null,
    token: null,
    identities: [],
    isLoading: true,
  });

  // Create AuthApi instance
  const authApi = useMemo(() => new AuthApi(), []);

  // Initialize auth state from localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem(TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        dispatch({ type: "SET_TOKEN", payload: storedToken });
        dispatch({ type: "SET_USER", payload: parsedUser });
      } catch (error) {
        console.error("Failed to parse stored user data:", error);
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }

    dispatch({ type: "SET_LOADING", payload: false });
  }, []);

  // Fetch identities when authenticated
  useEffect(() => {
    if (state.token) {
      refreshIdentities();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [state.token]);

  const refreshIdentities = async () => {
    if (!state.token) return;

    try {
      const data = await authApi.getIdentities();
      dispatch({ type: "SET_IDENTITIES", payload: data.identities || [] });
    } catch (error) {
      console.error("Failed to fetch identities:", error);
    }
  };

  const login = (newToken: string, newUser: User) => {
    localStorage.setItem(TOKEN_KEY, newToken);
    localStorage.setItem(USER_KEY, JSON.stringify(newUser));
    dispatch({ type: "LOGIN", payload: { token: newToken, user: newUser } });
  };

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    // Also remove OAuth token if it exists
    localStorage.removeItem("oauth_token");
    dispatch({ type: "LOGOUT" });
  };

  const linkProvider = async (provider: string) => {
    try {
      const data = await authApi.getProviderLoginUrl(provider);

      sessionStorage.setItem(`${provider}_auth_state`, data.state);
      sessionStorage.setItem("link_mode", "true");

      window.location.href = data.auth_url;
    } catch (error) {
      console.error(`Failed to initiate ${provider} login:`, error);
      throw error;
    }
  };

  const unlinkIdentity = async (identityId: string) => {
    if (!state.token) return;

    try {
      await authApi.unlinkIdentity(identityId);
      await refreshIdentities();
    } catch (error) {
      console.error("Failed to unlink identity:", error);
      throw error;
    }
  };

  const isAuthenticated = !!state.token && !!state.user;

  return (
    <AuthContext.Provider
      value={{
        user: state.user,
        token: state.token,
        identities: state.identities,
        isAuthenticated,
        isLoading: state.isLoading,
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
