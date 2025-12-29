import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRepositoryCreatePage } from "@/features/repository/hooks/useRepositoryCreatePage";
import { OAuthConnection } from "@/features/repository/components/OAuthConnection";
import { RepositorySelector } from "@/features/repository/components/RepositorySelector";
import { ConnectionStatus } from "@/features/repository/components/ConnectionStatus";
import { RepositoryConfirmation } from "@/features/repository/components/RepositoryConfirmation";
import {
  repositoryFormSchema,
  type RepositoryFormData,
} from "@/features/repository/types/repositoryForm";
import { LoadingSpinner, Alert, Card } from "@/ui";

/**
 * リポジトリ新規登録ページ
 * OAuth接続済みの場合はリポジトリ一覧から選択、未接続の場合はOAuth認証から開始
 */
export function RepositoryCreatePage() {
  const navigate = useNavigate();

  // ページ全体のロジックを統合するカスタムフック
  const {
    selectedProvider,
    selectedRepository,
    setSelectedRepository,
    isConnected,
    handleConnect,
    handleDisconnect,
    handleRegisterRepository,
    isLoadingConnections,
    isLoadingRepositories,
    isAuthenticating,
    isSubmitting,
    repositories,
    error,
    message,
  } = useRepositoryCreatePage();

  // フォーム管理
  const {
    register,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositoryFormSchema),
    defaultValues: {
      provider: "github",
    },
  });

  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  const handleCancel = () => {
    navigate("/repositories");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      {/* ローディング状態 */}
      {isLoadingConnections && (
        <Card>
          <LoadingSpinner message="Checking connections..." />
        </Card>
      )}

      {/* エラー表示 */}
      {error && <Alert type="error">{error}</Alert>}

      {!isLoadingConnections && (
        <>
          {/* Step 1: OAuth接続（未接続の場合） */}
          {!isConnected && (
            <OAuthConnection
              selectedProvider={selectedProvider}
              register={register}
              errors={errors}
              isAuthenticating={isAuthenticating}
              onConnect={handleConnect}
              gitlabUrl={gitlabUrl}
              gitlabClientId={gitlabClientId}
              gitlabClientSecret={gitlabClientSecret}
            />
          )}

          {/* 接続済みの場合: 接続状態を表示 */}
          {isConnected && (
            <ConnectionStatus
              provider={selectedProvider}
              onDisconnect={handleDisconnect}
            />
          )}

          {/* Step 2: リポジトリ選択（接続済みの場合） */}
          {isConnected && (
            <RepositorySelector
              repositories={repositories}
              isLoading={isLoadingRepositories}
              onSelect={setSelectedRepository}
              selectedRepository={selectedRepository}
            />
          )}

          {/* Step 3: 登録確認（リポジトリ選択済みの場合） */}
          {selectedRepository && (
            <RepositoryConfirmation
              repository={selectedRepository}
              isSubmitting={isSubmitting}
              message={message}
              onRegister={handleRegisterRepository}
              onCancel={handleCancel}
            />
          )}
        </>
      )}
    </div>
  );
}
