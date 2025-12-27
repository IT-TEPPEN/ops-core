import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

/**
 * OAuth2.0コールバックページ
 * GitプロバイダーからのOAuth認証後にリダイレクトされるページ
 */
function OAuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading"
  );
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get("code");
      const state = searchParams.get("state");
      const error = searchParams.get("error");
      const errorDescription = searchParams.get("error_description");

      // エラーがある場合
      if (error) {
        setStatus("error");
        setMessage(errorDescription || error || "OAuth authentication failed");
        setTimeout(() => {
          navigate("/repositories/new");
        }, 3000);
        return;
      }

      // codeとstateがない場合
      if (!code || !state) {
        setStatus("error");
        setMessage("Invalid OAuth callback - missing code or state");
        setTimeout(() => {
          navigate("/repositories/new");
        }, 3000);
        return;
      }

      try {
        // stateを検証
        const savedState = sessionStorage.getItem("oauth_state");
        const savedProvider = sessionStorage.getItem("oauth_provider");
        const savedGitlabUrl = sessionStorage.getItem("oauth_gitlab_url");
        const savedGitlabClientSecret = sessionStorage.getItem(
          "oauth_gitlab_client_secret"
        );

        if (!savedState || savedState !== state) {
          throw new Error("Invalid state parameter - possible CSRF attack");
        }

        if (!savedProvider) {
          throw new Error("OAuth provider not found in session");
        }

        // バックエンドにcodeを送信してトークンを取得
        const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
        const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

        const requestBody: {
          provider: string;
          code: string;
          state: string;
          gitlabUrl?: string;
          clientId?: string;
          clientSecret?: string;
        } = {
          provider: savedProvider,
          code,
          state,
        };

        // セルフホストGitLabの場合はURL、Client ID、Client Secretも送信
        if (savedProvider === "gitlab-self-hosted") {
          if (savedGitlabUrl) {
            requestBody.gitlabUrl = savedGitlabUrl;
          }
          if (savedGitlabClientSecret) {
            requestBody.clientSecret = savedGitlabClientSecret;
          }
          // Client IDもsessionStorageから取得して送信
          const savedGitlabClientId = sessionStorage.getItem(
            "oauth_gitlab_client_id"
          );
          if (savedGitlabClientId) {
            requestBody.clientId = savedGitlabClientId;
          }
        }

        const response = await fetch(`${apiUrl}/auth/oauth/callback`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody),
        });

        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.message || "Failed to exchange OAuth code");
        }

        // 認証成功
        setStatus("success");
        setMessage("Authentication successful! Redirecting...");

        // セッションストレージをクリーンアップ
        sessionStorage.removeItem("oauth_state");
        sessionStorage.removeItem("oauth_provider");
        sessionStorage.removeItem("oauth_gitlab_url");
        sessionStorage.removeItem("oauth_gitlab_client_id");
        sessionStorage.removeItem("oauth_gitlab_client_secret");

        // トークンをlocalStorageに保存（必要に応じて）
        if (data.access_token) {
          localStorage.setItem(
            `${savedProvider}_access_token`,
            data.access_token
          );
        }

        // リポジトリ登録ページにリダイレクト
        setTimeout(() => {
          navigate("/repositories/new", {
            state: { oauthSuccess: true, provider: savedProvider },
          });
        }, 2000);
      } catch (err) {
        setStatus("error");
        setMessage(
          err instanceof Error ? err.message : "An unknown error occurred"
        );
        setTimeout(() => {
          navigate("/repositories/new");
        }, 3000);
      }
    };

    handleCallback();
  }, [searchParams, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
      <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-lg max-w-md w-full">
        <div className="text-center">
          {status === "loading" && (
            <>
              <svg
                className="animate-spin h-12 w-12 mx-auto text-blue-600 mb-4"
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
              <h2 className="text-xl font-semibold mb-2">
                Processing Authentication...
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Please wait while we complete your authentication.
              </p>
            </>
          )}

          {status === "success" && (
            <>
              <svg
                className="h-12 w-12 mx-auto text-green-600 mb-4"
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
              <h2 className="text-xl font-semibold mb-2 text-green-600">
                Success!
              </h2>
              <p className="text-gray-600 dark:text-gray-400">{message}</p>
            </>
          )}

          {status === "error" && (
            <>
              <svg
                className="h-12 w-12 mx-auto text-red-600 mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <h2 className="text-xl font-semibold mb-2 text-red-600">
                Authentication Failed
              </h2>
              <p className="text-gray-600 dark:text-gray-400">{message}</p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default OAuthCallbackPage;
