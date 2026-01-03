export interface ValidateCodeAndGetTokenDto {
  get provider(): string;
  get code(): string;
  get state(): string;
  get rememberMe(): boolean;
}
