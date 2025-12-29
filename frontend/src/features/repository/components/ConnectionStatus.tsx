import { GitProvider } from "@/shared/api/gitProviderApi";
import { CheckCircleIcon, SpinnerIcon } from "@/ui";
import { useGitProvider } from "../hooks/useGitProvider";

interface ConnectionStatusProps {
  provider: GitProvider;
  onConnectionChange?: () => void;
}

export function ConnectionStatus({
  provider,
  onConnectionChange,
}: ConnectionStatusProps) {
  const { isConnected, disconnect, isLoadingConnections } = useGitProvider();
  const connected = isConnected(provider);

  const providerDisplayName = {
    github: "GitHub",
    gitlab: "GitLab",
    "gitlab-self-hosted": "Self-Hosted GitLab",
  }[provider];

  const handleDisconnect = async () => {
    await disconnect(provider);
    onConnectionChange?.();
  };

  if (isLoadingConnections) {
    return (
      <div className="bg-gray-50 dark:bg-gray-800 p-6 rounded-lg shadow">
        <div className="flex items-center gap-3">
          <SpinnerIcon className="text-blue-600" size={24} />
          <span className="text-gray-600 dark:text-gray-400">
            Checking connection status...
          </span>
        </div>
      </div>
    );
  }

  if (!connected) {
    return null;
  }

  return (
    <div className="bg-green-50 dark:bg-green-900/20 p-6 rounded-lg shadow border border-green-200 dark:border-green-800">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <CheckCircleIcon className="text-green-600" size={24} />
          <div>
            <h3 className="font-semibold text-green-800 dark:text-green-200">
              Connected to {providerDisplayName}
            </h3>
            <p className="text-sm text-green-600 dark:text-green-400">
              You can now select a repository from your account
            </p>
          </div>
        </div>
        <button
          type="button"
          onClick={handleDisconnect}
          className="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
        >
          Disconnect
        </button>
      </div>
    </div>
  );
}
