import { useEffect } from "react";
import {
  useNavigate,
  useSearchParams,
  useLocation,
  useParams,
} from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export default function AuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const { provider } = useParams<{ provider: string }>();
  const { login } = useAuth();

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

      if (state !== storedState || provider !== storedProvider) {
        console.error("State or provider mismatch - potential CSRF attack");
        navigate("/login", { replace: true });
        return;
      }

      try {
        // Exchange code for token
        const response = await fetch(
          `http://localhost:8080/api/v1/auth/${provider}/callback`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ code, state }),
          }
        );

        if (!response.ok) {
          throw new Error("Failed to authenticate");
        }

        const data = await response.json();

        // Store token and user info
        login(data.token, data.user);

        // Clean up
        sessionStorage.removeItem("auth_state");
        sessionStorage.removeItem("auth_provider");

        // Get the originally requested page or default to home
        const from = (location.state as any)?.from?.pathname || "/";
        navigate(from, { replace: true });
      } catch (error) {
        console.error("Authentication error:", error);
        navigate("/login", { replace: true });
      }
    };

    handleCallback();
  }, [searchParams, navigate, location, login, provider]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-4 text-gray-600">Completing sign in...</p>
      </div>
    </div>
  );
}
