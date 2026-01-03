export interface TokenDto {
  get accessToken(): string;
  get refreshToken(): string;
}
