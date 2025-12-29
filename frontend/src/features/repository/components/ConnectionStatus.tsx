import { GitProvider } from "@/shared/api/gitProviderApi";
import { CheckCircleIcon } from "@/ui";

interface ConnectionStatusProps {
  provider: GitProvider;
  onDisconnect: () => void;
}

export function ConnectionStatus({
  provider,
  onDisconnect,
}: ConnectionStatusProps) {
  const providerDisplayName = {
    github: "GitHub",
    gitlab: "GitLab",
    "gitlab-self-hosted": "Self-Hosted GitLab",
  }[provider];

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
          onClick={onDisconnect}
          className="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
        >
          Disconnect
        </button>
      </div>
    </div>
  );
}
