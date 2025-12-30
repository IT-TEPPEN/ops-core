import { useEffect, useMemo } from "react";
import {
  useNavigate,
  useSearchParams,
  useLocation,
  useParams,
} from "react-router-dom";
import { useAuth } from "../app/hooks/useAuth";
import { AuthApi } from "@/shared/api/authApi";

export default function AuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const { provider } = useParams<{ provider: string }>();
  const { login } = useAuth();
  const authApi = useMemo(() => new AuthApi(), []);

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get("code");
      const state = searchParams.get("state");

      if (!code || !state || !provider) {
        console.error("Missing code, state, or provider parameter");
        navigate("/login", { replace: true });
        return;
      }

      // Verify state for CSRF protection
      const storedState = sessionStorage.getItem("auth_state");
      const storedProvider = sessionStorage.getItem("auth_provider");
      const rememberMe = sessionStorage.getItem("auth_remember_me") === "true";

      if (state !== storedState || provider !== storedProvider) {
        console.error("State or provider mismatch - potential CSRF attack");
        navigate("/login", { replace: true });
        return;
      }

      try {
        // Exchange code for token (with remember_me preference)
        const data = await authApi.handleProviderCallback(
          provider,
          code,
          state,
          rememberMe
        );

        // Store token, refresh token, and user info
        login(data.token, data.user, data.refresh_token);

        // Clean up
        sessionStorage.removeItem("auth_state");
        sessionStorage.removeItem("auth_provider");
        sessionStorage.removeItem("auth_remember_me");

        // Get the originally requested page or default to home
        const from =
          (location.state as { from?: { pathname: string } } | null)?.from
            ?.pathname || "/";
        navigate(from, { replace: true });
      } catch (error) {
        console.error("Authentication error:", error);
        navigate("/login", { replace: true });
      }
    };

    handleCallback();
  }, [searchParams, navigate, location, login, provider, authApi]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-4 text-gray-600">Completing sign in...</p>
      </div>
    </div>
  );
}
