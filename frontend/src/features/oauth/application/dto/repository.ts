export interface Repository {
  get id(): number;
  get name(): string;
  get fullName(): string;
  get description(): string;
  get htmlUrl(): string;
  get cloneUrl(): string;
  get private(): boolean;
  get defaultBranch(): string;
  get owner(): {
    login: string;
    avatarUrl: string;
  };
}
