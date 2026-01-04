export interface StartLoginProcessDto {
  get provider(): string;
  get rememberMe(): boolean;
  get from(): string;
  redirectToAuthenticationPage: (externalUrl: string) => void;
}
