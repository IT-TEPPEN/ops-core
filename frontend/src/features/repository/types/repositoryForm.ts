import { GitProvider } from "@/shared/api/gitProviderApi";
import { z } from "zod";

export const gitProviders = ["github", "gitlab", "gitlab-self-hosted"] as const;

export const repositoryFormSchema = z
  .object({
    provider: z.enum(gitProviders),
    gitlabUrl: z.string().optional(),
    gitlabClientId: z.string().optional(),
    gitlabClientSecret: z.string().optional(),
    url: z.string().optional(),
  })
  .refine(
    (data) => {
      // gitlab-self-hostedの場合はgitlabUrl、gitlabClientId、gitlabClientSecretが必須
      if (data.provider === "gitlab-self-hosted") {
        return (
          !!data.gitlabUrl && !!data.gitlabClientId && !!data.gitlabClientSecret
        );
      }
      return true;
    },
    {
      message:
        "GitLab URL, Client ID, and Client Secret are required for self-hosted GitLab",
      path: ["gitlabUrl"],
    }
  );

export type RepositoryFormData = z.infer<typeof repositoryFormSchema>;

export function getProviderDisplayName(provider: GitProvider): string {
  const displayNames: Record<GitProvider, string> = {
    github: "GitHub",
    gitlab: "GitLab",
    "gitlab-self-hosted": "Self-Hosted GitLab",
  };
  return displayNames[provider];
}
