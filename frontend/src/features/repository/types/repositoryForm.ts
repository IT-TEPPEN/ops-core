import { GitProvider } from "@/shared/api/gitProviderApi";
import { z } from "zod";

export const gitProviders = ["github", "gitlab", "gitlab-self-hosted"] as const;

export const repositoryFormSchema = z.object({
  gitlabUrl: z.string().optional(),
  gitlabClientId: z.string().optional(),
  gitlabClientSecret: z.string().optional(),
  url: z.string().optional(),
});

export type RepositoryFormData = z.infer<typeof repositoryFormSchema>;

export function getProviderDisplayName(provider: GitProvider): string {
  const displayNames: Record<GitProvider, string> = {
    github: "GitHub",
    gitlab: "GitLab",
    "gitlab-self-hosted": "Self-Hosted GitLab",
  };
  return displayNames[provider];
}
