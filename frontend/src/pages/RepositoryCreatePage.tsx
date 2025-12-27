import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { UI_Form_Input, UI_Form_Submit } from "../ui/form";
import { UI_Form_Field } from "../components";
import {
  initiateOAuthFlow,
  type GitProvider,
  type SelfHostedOAuthParams,
} from "../utils/oauth";

// GitプロバイダーのタイプLiteral型定義
const gitProviders = ["github", "gitlab", "gitlab-self-hosted"] as const;

// Zodバリデーションスキーマ
const repositorySchema = z
  .object({
    provider: z.enum(gitProviders, {
      required_error: "Please select a Git provider",
    }),
    gitlabUrl: z.string().optional(),
    gitlabClientId: z.string().optional(),
    gitlabClientSecret: z.string().optional(),
    url: z
      .string()
      .min(1, "Repository URL is required")
      .url("Please enter a valid URL"),
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
 * 新しいリポジトリをシステムに登録する
 */
function RepositoryCreatePage() {
  const navigate = useNavigate();
  const [isAuthenticating, setIsAuthenticating] = useState(false);
  const [submitMessage, setSubmitMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositorySchema),
    defaultValues: {
      provider: "github",
    },
  });

  const selectedProvider = watch("provider");
  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  // API base URL - directly use the base URL to avoid recalculation
  const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
  const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

  const handleOAuthConnect = async () => {
    if (!selectedProvider) {
      setSubmitMessage({
        type: "error",
        text: "Please select a Git provider first",
      });
      return;
    }

    // セルフホストの場合の追加バリデーション
    if (selectedProvider === "gitlab-self-hosted") {
      if (!gitlabUrl || !gitlabClientId || !gitlabClientSecret) {
        setSubmitMessage({
          type: "error",
          text: "Please enter GitLab URL, Client ID, and Client Secret",
        });
        return;
      }
    }

    setIsAuthenticating(true);
    setSubmitMessage(null);

    try {
      // セルフホスト用のパラメータを準備
      const selfHostedParams: SelfHostedOAuthParams | undefined =
        selectedProvider === "gitlab-self-hosted"
          ? {
              gitlabUrl: gitlabUrl!,
              clientId: gitlabClientId!,
              clientSecret: gitlabClientSecret!,
            }
          : undefined;

      // OAuth認証フローを開始
      await initiateOAuthFlow(selectedProvider, selfHostedParams);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to initiate OAuth flow";
      setSubmitMessage({ type: "error", text: message });
      setIsAuthenticating(false);
    }
  };

  const onSubmit = async (data: RepositoryFormData) => {
    setSubmitMessage(null);

    try {
      const response = await fetch(`${apiUrl}/repositories`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          provider: data.provider,
          url: data.url,
        }),
      });

      const responseData = await response.json();

      if (!response.ok) {
        throw new Error(
          responseData.message || "Failed to register repository"
        );
      }

      setSubmitMessage({
        type: "success",
        text: "Repository registered successfully!",
      });
      reset();

      // 成功後、2秒待ってから一覧ページに遷移
      setTimeout(() => {
        navigate("/repositories");
      }, 2000);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      setSubmitMessage({ type: "error", text: message });
    }
  };

  const handleCancel = () => {
    navigate("/repositories");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      {/* OAuth接続セクション */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-lg font-semibold mb-4">
          Step 1: Connect Your Git Account
        </h2>
        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
          First, connect your Git account using OAuth2.0 to securely access your
          repositories.
        </p>

        <div className="space-y-4">
          <UI_Form_Field
            label="Git Provider"
            name="provider"
            error={errors.provider?.message}
            required
          >
            <select
              {...register("provider")}
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
                  Create an OAuth application in your GitLab instance at
                  Settings → Applications
                </p>
              </UI_Form_Field>
            </>
          )}

          <button
            type="button"
            onClick={handleOAuthConnect}
            disabled={isAuthenticating}
            className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isAuthenticating ? (
              <>
                <svg
                  className="animate-spin h-5 w-5"
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
                Connecting...
              </>
            ) : (
              <>
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
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

      {/* リポジトリ登録フォーム */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-lg font-semibold mb-4">
          Step 2: Register Repository
        </h2>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <UI_Form_Field
            label="Repository URL"
            name="repoUrl"
            error={errors.url?.message}
            required
          >
            <UI_Form_Input
              id="repoUrl"
              type="text"
              placeholder="https://github.com/username/repo.git"
              {...register("url")}
            />
          </UI_Form_Field>

          <div className="flex gap-3">
            <UI_Form_Submit
              isSubmitting={isSubmitting}
              label="Register Repository"
            />
            <button
              type="button"
              onClick={handleCancel}
              className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
              disabled={isSubmitting}
            >
              Cancel
            </button>
          </div>
        </form>

        {submitMessage && (
          <div
            className={`mt-4 p-3 rounded ${
              submitMessage.type === "success"
                ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
            }`}
          >
            {submitMessage.text}
          </div>
        )}
      </div>
    </div>
  );
}

export default RepositoryCreatePage;
