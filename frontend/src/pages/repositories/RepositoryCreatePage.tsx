import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRepositoryRegistration } from "@/features/repository/hooks/useRepositoryRegistration";
import { useGitProvider } from "@/features/repository/hooks/useGitProvider";
import { OAuthConnection } from "@/features/repository/components/OAuthConnection";
import { RepositorySelector } from "@/features/repository/components/RepositorySelector";
import { GitRepository, GitProvider } from "@/shared/api/gitProviderApi";

// GitプロバイダーのタイプLiteral型定義
const gitProviders = ["github", "gitlab", "gitlab-self-hosted"] as const;

// Zodバリデーションスキーマ
const repositorySchema = z
  .object({
    provider: z.enum(gitProviders),
    gitlabUrl: z.string().optional(),
    gitlabClientId: z.string().optional(),
    gitlabClientSecret: z.string().optional(),
    url: z.string().optional(),
  })
  .refine(
    (data) => {
      // gitlab-self-hostedの場合はgitlabUrl、gitlabClientId、gitlabClientSecretが必須
      if (data.provider === "gitlab-self-hosted") {
        return (
          !!data.gitlabUrl && !!data.gitlabClientId && !!data.gitlabClientSecret
        );
      }
      return true;
    },
    {
      message:
        "GitLab URL, Client ID, and Client Secret are required for self-hosted GitLab",
      path: ["gitlabUrl"],
    }
  );

type RepositoryFormData = z.infer<typeof repositorySchema>;

/**
 * リポジトリ新規登録ページ
 * OAuth接続済みの場合はリポジトリ一覧から選択、未接続の場合はOAuth認証から開始
 */
export function RepositoryCreatePage() {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    isAuthenticating,
    isSubmitting,
    message,
    handleOAuthConnect,
    handleSubmitRepository,
  } = useRepositoryRegistration();

  const {
    connections,
    repositories,
    isLoadingConnections,
    isLoadingRepositories,
    error: gitProviderError,
    isConnected,
    loadRepositories,
    disconnect,
    refreshConnections,
  } = useGitProvider();

  const [selectedRepository, setSelectedRepository] =
    useState<GitRepository | null>(null);

  const {
    register,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositorySchema),
    defaultValues: {
      provider: "github",
    },
  });

  const selectedProvider = watch("provider") as GitProvider;
  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  // OAuthコールバックからの遷移を検出
  const locationState = location.state as {
    oauthSuccess?: boolean;
    provider?: string;
  } | null;

  // OAuth成功後にリポジトリ一覧を読み込む
  useEffect(() => {
    if (locationState?.oauthSuccess && locationState?.provider) {
      refreshConnections();
    }
  }, [locationState, refreshConnections]);

  // プロバイダーが接続済みの場合、リポジトリ一覧を読み込む
  useEffect(() => {
    if (isConnected(selectedProvider)) {
      loadRepositories(selectedProvider);
    }
  }, [selectedProvider, isConnected, loadRepositories, connections]);

  const handleRepositorySelect = (repo: GitRepository) => {
    setSelectedRepository(repo);
  };

  const handleRegisterSelected = async () => {
    if (!selectedRepository) return;

    await handleSubmitRepository({
      provider: selectedProvider,
      url: selectedRepository.cloneUrl,
    });
  };

  const handleDisconnect = async () => {
    await disconnect(selectedProvider);
    setSelectedRepository(null);
  };

  const handleCancel = () => {
    navigate("/repositories");
  };

  const providerConnected = isConnected(selectedProvider);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      {/* ローディング状態 */}
      {isLoadingConnections && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <div className="flex items-center justify-center py-4">
            <svg
              className="animate-spin h-6 w-6 text-blue-600 mr-3"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              ></circle>
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <span>Checking connections...</span>
          </div>
        </div>
      )}

      {/* エラー表示 */}
      {gitProviderError && (
        <div className="bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200 p-4 rounded-lg">
          {gitProviderError}
        </div>
      )}

      {!isLoadingConnections && (
        <>
          {/* Step 1: OAuth接続（未接続の場合） */}
          {!providerConnected && (
            <OAuthConnection
              selectedProvider={selectedProvider}
              register={register}
              errors={errors}
              isAuthenticating={isAuthenticating}
              onConnect={handleOAuthConnect}
              gitlabUrl={gitlabUrl}
              gitlabClientId={gitlabClientId}
              gitlabClientSecret={gitlabClientSecret}
            />
          )}

          {/* 接続済みの場合: 接続状態を表示 */}
          {providerConnected && (
            <div className="bg-green-50 dark:bg-green-900/20 p-6 rounded-lg shadow border border-green-200 dark:border-green-800">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <svg
                    className="w-6 h-6 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <div>
                    <h3 className="font-semibold text-green-800 dark:text-green-200">
                      Connected to{" "}
                      {selectedProvider === "github"
                        ? "GitHub"
                        : selectedProvider === "gitlab"
                        ? "GitLab"
                        : "Self-Hosted GitLab"}
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
          )}

          {/* Step 2: リポジトリ選択（接続済みの場合） */}
          {providerConnected && (
            <RepositorySelector
              repositories={repositories}
              isLoading={isLoadingRepositories}
              onSelect={handleRepositorySelect}
              selectedRepository={selectedRepository}
            />
          )}

          {/* Step 3: 登録確認（リポジトリ選択済みの場合） */}
          {selectedRepository && (
            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
              <h3 className="text-lg font-semibold mb-4">
                Confirm Registration
              </h3>
              <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg mb-4">
                <div className="flex items-center gap-3">
                  <img
                    src={selectedRepository.owner.avatarUrl}
                    alt={selectedRepository.owner.login}
                    className="w-10 h-10 rounded-full"
                  />
                  <div>
                    <p className="font-medium">{selectedRepository.fullName}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {selectedRepository.cloneUrl}
                    </p>
                  </div>
                </div>
              </div>

              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={handleRegisterSelected}
                  disabled={isSubmitting}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  {isSubmitting ? (
                    <>
                      <svg
                        className="animate-spin h-4 w-4"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                        ></circle>
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                      </svg>
                      Registering...
                    </>
                  ) : (
                    "Register Repository"
                  )}
                </button>
                <button
                  type="button"
                  onClick={handleCancel}
                  className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                  disabled={isSubmitting}
                >
                  Cancel
                </button>
              </div>

              {message && (
                <div
                  className={`mt-4 p-3 rounded ${
                    message.type === "success"
                      ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                      : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
                  }`}
                >
                  {message.text}
                </div>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
