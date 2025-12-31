import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  OAuthConnection,
  ConnectionStatus,
  useGitProvider,
} from "@/features/oauth";
import {
  RepositorySelector,
  RepositoryConfirmation,
  useRepositoryRegistration,
} from "@/features/repository";
import {
  repositoryFormSchema,
  type RepositoryFormData,
} from "@/features/repository/types/repositoryForm";
import { GitProvider, GitRepository } from "@/shared/api/gitProviderApi";
import { Alert } from "@/ui";

/**
 * リポジトリ新規登録ページ
 * OAuth接続済みの場合はリポジトリ一覧から選択、未接続の場合はOAuth認証から開始
 */
export function RepositoryCreatePage() {
  const navigate = useNavigate();

  // ページレベルの状態管理
  const [selectedProvider, setSelectedProvider] =
    useState<GitProvider>("github");
  const [selectedRepository, setSelectedRepository] =
    useState<GitRepository | null>(null);

  // 各コンポーネントが自分で必要なロジックを呼び出す
  const { error: gitProviderError } = useGitProvider();
  const { isAuthenticating, handleOAuthConnect } = useRepositoryRegistration();

  // フォーム管理
  const {
    register,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositoryFormSchema),
  });

  const gitlabUrl = watch("gitlabUrl");
  const gitlabClientId = watch("gitlabClientId");
  const gitlabClientSecret = watch("gitlabClientSecret");

  const handleConnect = (
    provider: GitProvider,
    selfHostedParams?: {
      gitlabUrl: string;
      clientId: string;
      clientSecret: string;
    }
  ) => {
    handleOAuthConnect(provider, selfHostedParams);
  };

  const handleConnectionChange = () => {
    setSelectedRepository(null);
  };

  const handleCancel = () => {
    navigate("/repositories");
  };

  const handleSuccess = () => {
    navigate("/repositories");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      {/* エラー表示 */}
      {gitProviderError && <Alert type="error">{gitProviderError}</Alert>}

      {/* Step 1: OAuth接続 */}
      <OAuthConnection
        selectedProvider={selectedProvider}
        register={register}
        errors={errors}
        isAuthenticating={isAuthenticating}
        onConnect={handleConnect}
        gitlabUrl={gitlabUrl}
        gitlabClientId={gitlabClientId}
        gitlabClientSecret={gitlabClientSecret}
        onChangeSelectedProvider={setSelectedProvider}
      />

      {/* 接続状態表示 */}
      <ConnectionStatus
        provider={selectedProvider}
        onConnectionChange={handleConnectionChange}
      />

      {/* Step 2: リポジトリ選択 */}
      <RepositorySelector
        provider={selectedProvider}
        onSelect={setSelectedRepository}
        selectedRepository={selectedRepository}
      />

      {/* Step 3: 登録確認 */}
      {selectedRepository && (
        <RepositoryConfirmation
          repository={selectedRepository}
          provider={selectedProvider}
          onCancel={handleCancel}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  );
}
