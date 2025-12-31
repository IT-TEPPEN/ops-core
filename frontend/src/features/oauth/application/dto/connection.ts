export interface Connection {
  get id(): string;
  get provider(): string;
  get providerUsername(): string;
  get connectedAt(): Date;
}
