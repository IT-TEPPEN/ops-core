import { useForm } from "react-hook-form";
import { LoadingSpinner, UI_Form_Field, UI_Form_Input } from "@/ui";
import { LightningIcon } from "@/ui";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";

type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

const repositoryFormSchema = z.object({
  gitlabUrl: z.string().optional(),
  gitlabClientId: z.string().optional(),
  gitlabClientSecret: z.string().optional(),
  url: z.string().optional(),
});

export type RepositoryFormData = z.infer<typeof repositoryFormSchema>;

export function OAuthConnection() {
  const [selectedProvider, setSelectedProvider] =
    useState<GitProvider>("github");
  const { isAuthenticating, handleOAuthConnect } = useRepositoryRegistration();

  const {
    register,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositoryFormSchema),
  });

  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  const handleConnect = () => {
    // セルフホストの場合の追加バリデーション
    if (selectedProvider === "gitlab-self-hosted") {
      if (!gitlabUrl || !gitlabClientId || !gitlabClientSecret) {
        console.error(
          "[OAuthConnection] Missing self-hosted GitLab parameters"
        );
        return;
      }
      handleOAuthConnect(selectedProvider, {
        gitlabUrl,
        clientId: gitlabClientId,
        clientSecret: gitlabClientSecret,
      });
    } else {
      handleOAuthConnect(selectedProvider);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 p-6 max-w-4xl mx-auto my-4 rounded-lg shadow">
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
              setSelectedProvider(
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
            <LoadingSpinner message="Authenticating..." />
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

import { useReducer } from "react";

interface OAuthConfig {
  authUrl: string;
  clientId: string;
  scope: string;
  redirectUri: string;
}

interface SelfHostedOAuthParams {
  gitlabUrl: string;
  clientId: string;
  clientSecret: string;
}

/**
 * ランダムなstate文字列を生成（CSRF対策）
 */
function generateState(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return Array.from(array, (byte) => byte.toString(16).padStart(2, "0")).join(
    ""
  );
}

/**
 * プロバイダーごとのOAuth設定を取得
 * @param provider - GitプロバイダーType
 * @param selfHostedParams - セルフホスト用のパラメータ（gitlab-self-hostedの場合に必要）
 */
function getOAuthConfig(
  provider: GitProvider,
  selfHostedParams?: SelfHostedOAuthParams
): OAuthConfig {
  const redirectUri = `${window.location.origin}/oauth/callback`;

  switch (provider) {
    case "github":
      return {
        authUrl: "https://github.com/login/oauth/authorize",
        clientId: import.meta.env.VITE_GITHUB_CLIENT_ID || "",
        scope: "repo,read:user,user:email",
        redirectUri,
      };
    case "gitlab":
      return {
        authUrl: "https://gitlab.com/oauth/authorize",
        clientId: import.meta.env.VITE_GITLAB_CLIENT_ID || "",
        scope: "read_user read_api read_repository",
        redirectUri,
      };
    case "gitlab-self-hosted":
      if (!selfHostedParams?.gitlabUrl || !selfHostedParams?.clientId) {
        throw new Error(
          "GitLab URL and Client ID are required for self-hosted GitLab"
        );
      }
      // URLの末尾のスラッシュを削除
      const baseUrl = selfHostedParams.gitlabUrl.replace(/\/$/, "");
      return {
        authUrl: `${baseUrl}/oauth/authorize`,
        clientId: selfHostedParams.clientId,
        scope: "api,read_user,read_repository",
        redirectUri,
      };
    default:
      throw new Error(`Unsupported provider: ${provider}`);
  }
}

/**
 * OAuth認証フローを開始
 * @param provider - GitプロバイダーType ('github' | 'gitlab' | 'gitlab-self-hosted')
 * @param selfHostedParams - セルフホスト用のパラメータ（gitlab-self-hostedの場合に必要）
 */
export async function initiateOAuthFlow(
  provider: GitProvider,
  selfHostedParams?: SelfHostedOAuthParams
): Promise<void> {
  const config = getOAuthConfig(provider, selfHostedParams);

  if (!config.clientId) {
    const error = `OAuth client ID not configured for ${provider}. Please set VITE_${provider.toUpperCase()}_CLIENT_ID in your environment.`;
    throw new Error(error);
  }

  // CSRF対策のstateパラメータを生成
  const state = generateState();

  // stateとproviderをsessionStorageに保存（コールバック時に検証）
  sessionStorage.setItem("oauth_state", state);
  sessionStorage.setItem("oauth_provider", provider);

  // セルフホストの場合は追加情報も保存
  if (provider === "gitlab-self-hosted" && selfHostedParams) {
    sessionStorage.setItem("oauth_gitlab_url", selfHostedParams.gitlabUrl);
    sessionStorage.setItem("oauth_gitlab_client_id", selfHostedParams.clientId);
    sessionStorage.setItem(
      "oauth_gitlab_client_secret",
      selfHostedParams.clientSecret
    );
  }

  // OAuth認証URLを構築
  const params = new URLSearchParams({
    client_id: config.clientId,
    redirect_uri: config.redirectUri,
    scope: config.scope,
    state,
    response_type: "code",
  });

  const authUrl = `${config.authUrl}?${params.toString()}`;

  // OAuth認証ページにリダイレクト
  window.location.href = authUrl;
}

/**
 * GitProviderの型をエクスポート（他のファイルで使用するため）
 */
export type { SelfHostedOAuthParams };

type RegistrationState = {
  isAuthenticating: boolean;
  isSubmitting: boolean;
  message: {
    type: "success" | "error";
    text: string;
  } | null;
};

type RegistrationAction =
  | { type: "START_AUTH" }
  | { type: "AUTH_ERROR"; error: string }
  | { type: "START_SUBMIT" }
  | { type: "SUBMIT_SUCCESS"; message: string }
  | { type: "SUBMIT_ERROR"; error: string }
  | { type: "CLEAR_MESSAGE" };

const initialState: RegistrationState = {
  isAuthenticating: false,
  isSubmitting: false,
  message: null,
};

function registrationReducer(
  state: RegistrationState,
  action: RegistrationAction
): RegistrationState {
  switch (action.type) {
    case "START_AUTH":
      return { ...state, isAuthenticating: true, message: null };
    case "AUTH_ERROR":
      return {
        ...state,
        isAuthenticating: false,
        message: { type: "error", text: action.error },
      };
    case "START_SUBMIT":
      return { ...state, isSubmitting: true, message: null };
    case "SUBMIT_SUCCESS":
      return {
        ...state,
        isSubmitting: false,
        message: { type: "success", text: action.message },
      };
    case "SUBMIT_ERROR":
      return {
        ...state,
        isSubmitting: false,
        message: { type: "error", text: action.error },
      };
    case "CLEAR_MESSAGE":
      return { ...state, message: null };
    default:
      return state;
  }
}

export function useRepositoryRegistration() {
  const [state, dispatch] = useReducer(registrationReducer, initialState);

  const handleOAuthConnect = async (
    provider: "github" | "gitlab" | "gitlab-self-hosted",
    selfHostedParams?: SelfHostedOAuthParams
  ) => {
    dispatch({ type: "START_AUTH" });

    try {
      await initiateOAuthFlow(provider, selfHostedParams);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to initiate OAuth flow";
      dispatch({ type: "AUTH_ERROR", error: message });
    }
  };

  const clearMessage = () => {
    dispatch({ type: "CLEAR_MESSAGE" });
  };

  return {
    ...state,
    handleOAuthConnect,
    clearMessage,
  };
}
