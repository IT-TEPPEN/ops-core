import { UseFormRegister, FieldErrors } from "react-hook-form";
import { SelfHostedOAuthParams } from "@/shared/utils/oauth";
import { UI_Form_Field, UI_Form_Input } from "@/ui";
import { SpinnerIcon, LightningIcon } from "@/ui";

type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

interface OAuthConnectionFormData {
  url?: string;
  gitlabUrl?: string;
  gitlabClientId?: string;
  gitlabClientSecret?: string;
}

interface OAuthConnectionProps {
  selectedProvider: GitProvider;
  register: UseFormRegister<OAuthConnectionFormData>;
  errors: FieldErrors<OAuthConnectionFormData>;
  isAuthenticating: boolean;
  onConnect: (
    provider: GitProvider,
    selfHostedParams?: SelfHostedOAuthParams
  ) => void;
  gitlabUrl?: string;
  gitlabClientId?: string;
  gitlabClientSecret?: string;
  onChangeSelectedProvider: (provider: GitProvider) => void;
}

export function OAuthConnection({
  selectedProvider,
  register,
  errors,
  isAuthenticating,
  onConnect,
  gitlabUrl,
  gitlabClientId,
  gitlabClientSecret,
  onChangeSelectedProvider,
}: OAuthConnectionProps) {
  const handleConnect = () => {
    console.log(
      "[OAuthConnection] handleConnect called with provider:",
      selectedProvider
    );
    // セルフホストの場合の追加バリデーション
    if (selectedProvider === "gitlab-self-hosted") {
      if (!gitlabUrl || !gitlabClientId || !gitlabClientSecret) {
        console.error(
          "[OAuthConnection] Missing self-hosted GitLab parameters"
        );
        return;
      }
      console.log(
        "[OAuthConnection] Calling onConnect with self-hosted params"
      );
      onConnect(selectedProvider, {
        gitlabUrl,
        clientId: gitlabClientId,
        clientSecret: gitlabClientSecret,
      });
    } else {
      console.log(
        "[OAuthConnection] Calling onConnect for provider:",
        selectedProvider
      );
      onConnect(selectedProvider);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">
        Step 1: Connect Your Git Account
      </h2>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        First, connect your Git account using OAuth2.0 to securely access your
        repositories.
      </p>

      <div className="space-y-4">
        <UI_Form_Field label="Git Provider" name="provider" error={""} required>
          <select
            onChange={(e) => {
              e.preventDefault();
              onChangeSelectedProvider(
                e.currentTarget.value as unknown as GitProvider
              );
            }}
            value={selectedProvider}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="github">GitHub</option>
            <option value="gitlab">GitLab (gitlab.com)</option>
            <option value="gitlab-self-hosted">
              GitLab (Self-Hosted / Local)
            </option>
          </select>
        </UI_Form_Field>

        {/* セルフホストGitLabの場合の追加フィールド */}
        {selectedProvider === "gitlab-self-hosted" && (
          <>
            <UI_Form_Field
              label="GitLab URL"
              name="gitlabUrl"
              error={errors.gitlabUrl?.message}
              required
            >
              <UI_Form_Input
                id="gitlabUrl"
                type="text"
                placeholder="https://gitlab.example.com"
                {...register("gitlabUrl")}
              />
            </UI_Form_Field>

            <UI_Form_Field
              label="Client ID"
              name="gitlabClientId"
              error={errors.gitlabClientId?.message}
              required
            >
              <UI_Form_Input
                id="gitlabClientId"
                type="text"
                placeholder="Your GitLab OAuth Application Client ID"
                {...register("gitlabClientId")}
              />
            </UI_Form_Field>

            <UI_Form_Field
              label="Client Secret"
              name="gitlabClientSecret"
              error={errors.gitlabClientSecret?.message}
              required
            >
              <UI_Form_Input
                id="gitlabClientSecret"
                type="password"
                placeholder="Your GitLab OAuth Application Client Secret"
                {...register("gitlabClientSecret")}
              />
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Create an OAuth application in your GitLab instance at Settings
                → Applications
              </p>
            </UI_Form_Field>
          </>
        )}

        <button
          type="button"
          onClick={handleConnect}
          disabled={isAuthenticating}
          className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {isAuthenticating ? (
            <>
              <SpinnerIcon size={20} />
              Connecting...
            </>
          ) : (
            <>
              <LightningIcon size={20} />
              Connect with{" "}
              {selectedProvider === "github"
                ? "GitHub"
                : selectedProvider === "gitlab"
                ? "GitLab"
                : "Self-Hosted GitLab"}
            </>
          )}
        </button>
      </div>
    </div>
  );
}
