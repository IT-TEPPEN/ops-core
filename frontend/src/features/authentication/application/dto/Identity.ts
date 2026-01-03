export interface IdentityDto {
  get id(): string;
  get provider(): string;
  get email(): string;
  get name(): string;
  get pictureUrl(): string;
  get linkedAt(): string;
  get lastUsedAt(): string;
}

export interface IdentitiesDto {
  get identities(): IdentityDto[];
}
