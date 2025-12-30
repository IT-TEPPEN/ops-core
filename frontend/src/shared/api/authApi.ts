/**
 * 認証関連APIクライアント
 */
import { V1ApiClient } from "./client";

// --- Response Types ---

export interface ProviderLoginResponse {
  auth_url: string;
  state: string;
}

export interface AuthCallbackResponse {
  token: string;
  refresh_token: string;
  user: UserProfile;
}

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  picture: string;
}

export interface Identity {
  id: string;
  provider: string;
  email: string;
  name: string;
  linkedAt: string;
  lastUsedAt: string;
}

export interface IdentitiesResponse {
  identities: Identity[];
}

export interface OAuthConnectRequest {
  provider: string;
  code: string;
  state: string;
  gitlabUrl?: string;
  clientId?: string;
  clientSecret?: string;
}

export interface OAuthConnectionResponse {
  id: string;
  provider: string;
  providerUsername: string;
  scopes: string[];
  connectedAt: string;
}

// --- Auth API Adapter Interface ---

export interface AuthApiAdapter {
  /**
   * プロバイダーログインのURLを取得
   */
  getProviderLoginUrl(provider: string): Promise<ProviderLoginResponse>;

  /**
   * プロバイダーコールバックを処理してJWTトークンを取得
   */
  handleProviderCallback(
    provider: string,
    code: string,
    state: string
  ): Promise<AuthCallbackResponse>;

  /**
   * ユーザーのアイデンティティ一覧を取得
   */
  getIdentities(): Promise<IdentitiesResponse>;

  /**
   * アイデンティティを削除（リンク解除）
   */
  unlinkIdentity(identityId: string): Promise<void>;

  /**
   * OAuth接続を保存（Git プロバイダー接続用）
   */
  connectOAuth(request: OAuthConnectRequest): Promise<OAuthConnectionResponse>;

  /**
   * ログアウト
   */
  logout(): Promise<void>;
}

// --- Auth API Implementation ---

export class AuthApi extends V1ApiClient implements AuthApiAdapter {
  constructor() {
    super("/auth");
  }

  async getProviderLoginUrl(provider: string): Promise<ProviderLoginResponse> {
    return this.get<ProviderLoginResponse>(`/${provider}/login`);
  }

  async handleProviderCallback(
    provider: string,
    code: string,
    state: string,
    rememberMe: boolean = false
  ): Promise<AuthCallbackResponse> {
    return this.post<AuthCallbackResponse, { code: string; state: string; remember_me: boolean }>(
      `/${provider}/callback`,
      { code, state, remember_me: rememberMe }
    );
  }

  async refreshToken(refreshToken: string): Promise<AuthCallbackResponse> {
    return this.post<AuthCallbackResponse, { refresh_token: string }>(
      `/refresh`,
      { refresh_token: refreshToken }
    );
  }

  async getIdentities(): Promise<IdentitiesResponse> {
    return this.get<IdentitiesResponse>("/identities");
  }

  async unlinkIdentity(identityId: string): Promise<void> {
    await this.delete<void>(`/identities/${identityId}`);
  }

  async connectOAuth(
    request: OAuthConnectRequest
  ): Promise<OAuthConnectionResponse> {
    return this.post<OAuthConnectionResponse, OAuthConnectRequest>(
      "/oauth/connect",
      request
    );
  }

  async logout(): Promise<void> {
    await this.post<void, Record<string, never>>("/logout", {});
  }
}
