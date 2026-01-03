export interface StartLoginProcessDto {
  get provider(): string;
  get rememberMe(): boolean;
  redirectToAuthenticationPage: (externalUrl: string) => void;
}
