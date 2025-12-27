/**
 * OAuth2.0ヘルパー関数
 * GitプロバイダーとのOAuth認証フローを管理
 */

type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

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
        scope: "api,read_user,read_repository",
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
    throw new Error(
      `OAuth client ID not configured for ${provider}. Please set VITE_${provider.toUpperCase()}_CLIENT_ID in your environment.`
    );
  }

  // CSRF対策のstateパラメータを生成
  const state = generateState();

  // stateとproviderをsessionStorageに保存（コールバック時に検証）
  sessionStorage.setItem("oauth_state", state);
  sessionStorage.setItem("oauth_provider", provider);

  // セルフホストの場合は追加情報も保存
  if (provider === "gitlab-self-hosted" && selfHostedParams) {
    sessionStorage.setItem("oauth_gitlab_url", selfHostedParams.gitlabUrl);
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
 * OAuth認証が完了しているかチェック
 * @param provider - GitプロバイダーType
 */
export function isOAuthAuthenticated(provider: GitProvider): boolean {
  const token = localStorage.getItem(`${provider}_access_token`);
  return !!token;
}

/**
 * 保存されているOAuthトークンを取得
 * @param provider - GitプロバイダーType
 */
export function getOAuthToken(provider: GitProvider): string | null {
  return localStorage.getItem(`${provider}_access_token`);
}

/**
 * OAuthトークンを削除（ログアウト）
 * @param provider - GitプロバイダーType
 */
export function removeOAuthToken(provider: GitProvider): void {
  localStorage.removeItem(`${provider}_access_token`);
}

/**
 * すべてのプロバイダーのOAuthトークンを削除
 */
export function removeAllOAuthTokens(): void {
  removeOAuthToken("github");
  removeOAuthToken("gitlab");
  removeOAuthToken("gitlab-self-hosted");
}

/**
 * GitProviderの型をエクスポート（他のファイルで使用するため）
 */
export type { GitProvider, SelfHostedOAuthParams };
