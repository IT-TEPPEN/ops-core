import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  OAuthConnection,
  ConnectionList,
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
import { Connection } from "@/features/oauth/application/dto";

/**
 * リポジトリ新規登録ページ
 * OAuth接続済みの場合はリポジトリ一覧から選択、未接続の場合はOAuth認証から開始
 */
export function RepositoryCreatePage() {
  const navigate = useNavigate();
  const previousConnectionCountRef = useRef(0);

  // ページレベルの状態管理
  const [mode, setMode] = useState<"add" | "select">("select");
  const [selectedConnection, setSelectedConnection] =
    useState<Connection | null>(null);
  const [selectedRepository, setSelectedRepository] =
    useState<GitRepository | null>(null);

  // 各コンポーネントが自分で必要なロジックを呼び出す
  const {
    connections,
    isLoadingConnections,
    error: gitProviderError,
  } = useGitProvider();
  const { isAuthenticating, handleOAuthConnect } = useRepositoryRegistration();

  // コネクション選択ハンドラー
  const handleSelectConnection = (connection: Connection) => {
    setSelectedConnection(connection);
    setSelectedRepository(null); // リポジトリ選択をリセット
    setMode("select"); // 選択モードに切り替え
  };

  // コネクション追加モードに切り替え
  const handleAddConnection = () => {
    setMode("add");
    setSelectedConnection(null);
    setSelectedRepository(null);
  };

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
    setSelectedConnection(null); // コネクション選択をリセット
  };

  const handleCancel = () => {
    navigate("/repositories");
  };

  const handleSuccess = () => {
    navigate("/repositories");
  };

  // 初回ロード時にコネクションがある場合は最初のものを選択
  useEffect(() => {
    if (connections.length > 0 && !selectedConnection && mode === "select") {
      setSelectedConnection(connections[0]);
    }
  }, [connections, selectedConnection, mode]);

  // コネクション数が増えたら（新規追加された可能性）selectモードに戻る
  useEffect(() => {
    if (
      mode === "add" &&
      connections.length > previousConnectionCountRef.current &&
      connections.length > 0
    ) {
      // 新しく追加されたコネクションを選択
      const newConnection = connections[connections.length - 1];
      setSelectedConnection(newConnection);
      setMode("select");
    }
    previousConnectionCountRef.current = connections.length;
  }, [connections, mode]);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Register New Repository</h1>

      {/* エラー表示 */}
      {gitProviderError && <Alert type="error">{gitProviderError}</Alert>}

      <div className="grid grid-cols-12 gap-6">
        {/* 左ペイン: コネクション一覧 */}
        <div className="col-span-12 lg:col-span-4 xl:col-span-3">
          <div className="sticky top-6">
            <ConnectionList
              connections={connections}
              selectedConnectionId={selectedConnection?.id || null}
              onSelectConnection={handleSelectConnection}
              onAddConnection={handleAddConnection}
              isLoading={isLoadingConnections}
            />
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
                  setSelectedConnection(null);
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
                    provider={selectedConnection.provider as GitProvider}
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
