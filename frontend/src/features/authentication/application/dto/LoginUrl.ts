export interface LoginUriDto {
  get authUrl(): string;
  get state(): string;
}
