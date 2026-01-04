import { useReducer, useEffect, ReactNode } from "react";
import { AuthContext } from "./AuthContextDefinition";

interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

const TOKEN_KEY = "auth_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const USER_KEY = "auth_user";

interface AuthProviderProps {
  children: ReactNode;
}

// State type
interface AuthState {
  token: string | null;
  refreshToken: string | null;
  isLoading: boolean;
}

// Action types
type AuthAction =
  | { type: "SET_TOKEN"; payload: string | null }
  | { type: "SET_REFRESH_TOKEN"; payload: string | null }
  | { type: "SET_USER"; payload: User | null }
  | { type: "SET_LOADING"; payload: boolean }
  | {
      type: "LOGIN";
      payload: { token: string; refreshToken: string; user: User };
    };

// Reducer
function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case "SET_TOKEN":
      return { ...state, token: action.payload };
    case "SET_REFRESH_TOKEN":
      return { ...state, refreshToken: action.payload };
    case "SET_LOADING":
      return { ...state, isLoading: action.payload };
    case "LOGIN":
      return {
        ...state,
        token: action.payload.token,
        refreshToken: action.payload.refreshToken,
      };
    default:
      return state;
  }
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [state, dispatch] = useReducer(authReducer, {
    token: null,
    refreshToken: null,
    isLoading: true,
  });

  // Initialize auth state from localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem(TOKEN_KEY);
    const storedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedToken && storedRefreshToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        dispatch({ type: "SET_TOKEN", payload: storedToken });
        dispatch({ type: "SET_REFRESH_TOKEN", payload: storedRefreshToken });
        dispatch({ type: "SET_USER", payload: parsedUser });
      } catch (error) {
        console.error("Failed to parse stored user data:", error);
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }

    dispatch({ type: "SET_LOADING", payload: false });
  }, []);

  const login = (newToken: string, newUser: User, newRefreshToken: string) => {
    localStorage.setItem(TOKEN_KEY, newToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken);
    localStorage.setItem(USER_KEY, JSON.stringify(newUser));
    dispatch({
      type: "LOGIN",
      payload: {
        token: newToken,
        refreshToken: newRefreshToken,
        user: newUser,
      },
    });
  };

  const isAuthenticated = !!state.token;

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading: state.isLoading,
        login,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
