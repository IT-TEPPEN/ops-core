import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRepositoryRegistration } from "../features/repository/hooks/useRepositoryRegistration";
import { OAuthConnection } from "../features/repository/components/OAuthConnection";
import { RepositoryRegistrationForm } from "../features/repository/components/RepositoryRegistrationForm";

// GitプロバイダーのタイプLiteral型定義
const gitProviders = ["github", "gitlab", "gitlab-self-hosted"] as const;

// Zodバリデーションスキーマ
const repositorySchema = z
  .object({
    provider: z.enum(gitProviders),
    gitlabUrl: z.string().optional(),
    gitlabClientId: z.string().optional(),
    gitlabClientSecret: z.string().optional(),
    url: z
      .string()
      .min(1, "Repository URL is required")
      .url("Please enter a valid URL"),
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

type RepositoryFormData = z.infer<typeof repositorySchema>;

/**
 * リポジトリ新規登録ページ
 * 新しいリポジトリをシステムに登録する
 */
function RepositoryCreatePage() {
  const navigate = useNavigate();
  const {
    isAuthenticating,
    isSubmitting,
    message,
    handleOAuthConnect,
    handleSubmitRepository,
  } = useRepositoryRegistration();

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositorySchema),
    defaultValues: {
      provider: "github",
    },
  });

  const selectedProvider = watch("provider");
  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  const onSubmit = async (data: RepositoryFormData) => {
    await handleSubmitRepository({
      provider: data.provider,
      url: data.url,
    });
  };

  const handleCancel = () => {
    navigate("/repositories");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      <OAuthConnection
        selectedProvider={selectedProvider}
        register={register}
        errors={errors}
        isAuthenticating={isAuthenticating}
        onConnect={handleOAuthConnect}
        gitlabUrl={gitlabUrl}
        gitlabClientId={gitlabClientId}
        gitlabClientSecret={gitlabClientSecret}
      />

      <form onSubmit={handleSubmit(onSubmit)}>
        <RepositoryRegistrationForm
          register={register}
          errors={errors}
          isSubmitting={isSubmitting}
          onCancel={handleCancel}
          message={message}
        />
      </form>
    </div>
  );
}

export default RepositoryCreatePage;
