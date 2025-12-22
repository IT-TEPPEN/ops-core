export interface Repository {
  getId(): string;
  getName(): string;
  getUrl(): string;
  getCreatedAt(): Date;
  getUpdatedAt(): Date;
}
