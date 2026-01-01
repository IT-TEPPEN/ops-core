import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  OAuthConnection,
  ConnectionList,
  useOAuthQueryService,
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
import { useQuery } from "@tanstack/react-query";

/**
 * リポジトリ新規登録ページ
 * OAuth接続済みの場合はリポジトリ一覧から選択、未接続の場合はOAuth認証から開始
 */
export function RepositoryCreatePage() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const selected_connection_id = searchParams.get("selected_connection_id");
  const mode = selected_connection_id ? "select" : "add";
  const oauthQueryService = useOAuthQueryService();
  const query = useQuery({
    queryKey: ["connections"],
    queryFn: async () => oauthQueryService.listConnections(),
  });
  const [selectedRepository, setSelectedRepository] =
    useState<GitRepository | null>(null);

  // 各コンポーネントが自分で必要なロジックを呼び出す
  const { isAuthenticating, handleOAuthConnect } = useRepositoryRegistration();

  // フォーム管理
  const {
    register,
    watch,
    formState: { errors },
  } = useForm<RepositoryFormData>({
    resolver: zodResolver(repositoryFormSchema),
  });

  if (query.isLoading) {
    return <div className="text-sm text-gray-500">Loading connections...</div>;
  }

  if (query.isError || !query.data) {
    return (
      <div className="text-sm text-red-500">
        Failed to load connections. Please try again.
      </div>
    );
  }

  const connections = query.data;
  const selectedConnection =
    connections.find(
      (conn, i) =>
        conn.id === selected_connection_id ||
        (selected_connection_id === "first" && i === 0)
    ) || null;

  // コネクション追加モードに切り替え
  const handleAddConnection = () => {
    setSearchParams((searchParams) => {
      searchParams.delete("selected_connection_id");
      return searchParams;
    });
  };

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
    setSearchParams((searchParams) => {
      searchParams.set("selected_connection_id", "first");
      return searchParams;
    });
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

      <div className="grid grid-cols-12 gap-6">
        {/* 左ペイン: コネクション一覧 */}
        <div className="col-span-12 lg:col-span-4 xl:col-span-3">
          <div className="sticky top-6">
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold">Connections</h2>
                <button
                  onClick={handleAddConnection}
                  className="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                >
                  + Add
                </button>
              </div>
              <ConnectionList />
            </div>
          </div>
        </div>

        {/* 右ペイン: リポジトリ登録フロー */}
        <div className="col-span-12 lg:col-span-8 xl:col-span-9 space-y-6">
          {mode === "add" ? (
            // モード1: コネクション追加
            <>
              <OAuthConnection
                selectedProvider={"github"}
                register={register}
                errors={errors}
                isAuthenticating={isAuthenticating}
                onConnect={handleConnect}
                gitlabUrl={gitlabUrl}
                gitlabClientId={gitlabClientId}
                gitlabClientSecret={gitlabClientSecret}
                onChangeSelectedProvider={() => {
                  // プロバイダー変更時、選択されたコネクションをリセット
                  setSearchParams((searchParams) => {
                    searchParams.set("selected_connection_id", "first");
                    return searchParams;
                  });
                }}
              />
            </>
          ) : (
            // モード2: リポジトリ選択
            <>
              {!selectedConnection ? (
                <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
                  <p className="text-gray-600 dark:text-gray-400">
                    Please select a connection from the left pane to view
                    repositories, or click &quot;Add&quot; to connect a new Git
                    provider.
                  </p>
                </div>
              ) : (
                <>
                  {/* Step 2: リポジトリ選択 */}
                  <RepositorySelector
                    onSelect={setSelectedRepository}
                    selectedRepository={selectedRepository}
                    selectedConnection={selectedConnection}
                  />

                  {/* Step 3: 登録確認 */}
                  {selectedRepository && (
                    <RepositoryConfirmation
                      repository={selectedRepository}
                      provider={selectedConnection.provider as GitProvider}
                      onCancel={handleCancel}
                      onSuccess={handleSuccess}
                    />
                  )}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
