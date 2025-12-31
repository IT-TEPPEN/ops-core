export type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

export interface GitRepository {
  id: number;
  name: string;
  fullName: string;
  description: string;
  htmlUrl: string;
  cloneUrl: string;
  private: boolean;
  defaultBranch: string;
  owner: {
    login: string;
    avatarUrl: string;
  };
}
