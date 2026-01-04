import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  GoogleIcon,
  GitHubIcon,
  GitLabIcon,
  MicrosoftIcon,
  LoadingSpinner,
} from "@/ui";
import {
  useGetUserIdentityUsecase,
  useStartLoginProcessUsecase,
} from "@/features/authentication/presentation/contexts";
import { useQuery } from "@tanstack/react-query";

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
  const [error, setError] = useState<string | null>(null);
  const [rememberMe, setRememberMe] = useState(false);
  const startLoginProcessUsecase = useStartLoginProcessUsecase();
  const getUserIdentityUsecase = useGetUserIdentityUsecase();
  const query = useQuery({
    queryKey: ["UserIdentity"],
    queryFn: async () => getUserIdentityUsecase.execute(),
  });

  useEffect(() => {
    if (query.data) {
      navigate("/", { replace: true });
    }
  }, [query.data, navigate]);

  if (query.isLoading) {
    return <LoadingSpinner message="loading..." />;
  }

  if (query.isError) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          Error loading user identity: {`${query.error}`}
        </div>
      </div>
    );
  }

  if (query.data) {
    return null;
  }

  const handleProviderLogin = async (provider: Provider) => {
    try {
      setError(null);
      await startLoginProcessUsecase.execute({
        provider,
        rememberMe,
        redirectToAuthenticationPage: (externalUrl: string) => {
          window.location.href = externalUrl;
        },
      });
    } catch (error) {
      console.error(`Failed to initiate ${provider} login:`, error);
      setError(`Failed to initiate ${provider} login. Please try again.`);
    }
  };

  return (
    <div className="h-full flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
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

          {/* Remember Me checkbox */}
          <div className="flex items-center justify-center mb-4">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              checked={rememberMe}
              onChange={(e) => setRememberMe(e.target.checked)}
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label
              htmlFor="remember-me"
              className="ml-2 block text-sm text-gray-700"
            >
              Remember me for 30 days (otherwise 7 days)
            </label>
          </div>

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
