export interface Connection {
  get id(): string;
  get provider(): string;
  get providerHost(): string;
  get providerUsername(): string;
  get connectedAt(): Date;
}
