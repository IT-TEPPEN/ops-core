export interface OAuthCommandService {
  initiateOAuthFlow(provider: string, redirectUri: string): Promise<string>;
}
