import { useEffect, useState, useMemo } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "../app/hooks/useAuth";
import { AuthApi } from "@/shared/api/authApi";
import { GoogleIcon, GitHubIcon, GitLabIcon, MicrosoftIcon } from "@/ui";

type Provider = "google" | "github" | "gitlab" | "microsoft";

interface ProviderInfo {
  name: string;
  icon: React.ReactElement;
  bgColor: string;
  hoverBgColor: string;
}

const providerInfo: Record<Provider, ProviderInfo> = {
  google: {
    name: "Google",
    icon: <GoogleIcon className="mr-2 size-5" />,
    bgColor: "bg-white",
    hoverBgColor: "hover:bg-gray-50",
  },
  github: {
    name: "GitHub",
    icon: <GitHubIcon className="mr-2 size-5" />,
    bgColor: "bg-gray-900 text-white",
    hoverBgColor: "hover:bg-gray-800",
  },
  gitlab: {
    name: "GitLab",
    icon: <GitLabIcon className="mr-2 size-5" />,
    bgColor: "bg-white",
    hoverBgColor: "hover:bg-gray-50",
  },
  microsoft: {
    name: "Microsoft",
    icon: <MicrosoftIcon className="mr-2 size-5" />,
    bgColor: "bg-white",
    hoverBgColor: "hover:bg-gray-50",
  },
};

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, isLoading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const authApi = useMemo(() => new AuthApi(), []);

  const from =
    (location.state as { from?: { pathname: string } } | null)?.from
      ?.pathname || "/";

  // Redirect if already authenticated
  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate, from]);

  const handleProviderLogin = async (provider: Provider) => {
    try {
      setError(null);
      const data = await authApi.getProviderLoginUrl(provider);

      // Store state and provider for CSRF protection
      sessionStorage.setItem("auth_state", data.state);
      sessionStorage.setItem("auth_provider", provider);

      // Redirect to provider's authorization page
      window.location.href = data.auth_url;
    } catch (error) {
      console.error(`Failed to initiate ${provider} login:`, error);
      setError(`Failed to initiate ${provider} login. Please try again.`);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to OpsCore
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Manage your operations documentation and procedures
          </p>
        </div>
        <div className="mt-8 space-y-4">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative">
              <span className="block sm:inline">{error}</span>
            </div>
          )}
          {(["google", "github", "gitlab", "microsoft"] as Provider[]).map(
            (provider) => {
              const info = providerInfo[provider];
              return (
                <button
                  key={provider}
                  onClick={() => handleProviderLogin(provider)}
                  className={`group relative w-full flex justify-center items-center py-3 px-4 border border-gray-300 text-sm font-medium rounded-md ${info.bgColor} ${info.hoverBgColor} focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors`}
                >
                  {info.icon}
                  Sign in with {info.name}
                </button>
              );
            }
          )}
        </div>
      </div>
    </div>
  );
}
