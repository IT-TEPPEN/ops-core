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

export interface OAuthCommandService {
  initiateOAuthFlow(provider: string, redirectUri: string): Promise<string>;

  connectOAuth(request: OAuthConnectRequest): Promise<OAuthConnectionResponse>;
}
