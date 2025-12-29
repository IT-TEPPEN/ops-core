import { useState, useEffect, useCallback } from "react";
import { useLocation } from "react-router-dom";
import { useRepositoryRegistration } from "./useRepositoryRegistration";
import { useGitProvider } from "./useGitProvider";
import { GitRepository, GitProvider } from "@/shared/api/gitProviderApi";

export interface RepositoryCreatePageState {
  // Provider selection
  selectedProvider: GitProvider;
  setSelectedProvider: (provider: GitProvider) => void;

  // Repository selection
  selectedRepository: GitRepository | null;
  setSelectedRepository: (repo: GitRepository | null) => void;

  // Connection management
  isConnected: boolean;
  handleConnect: (
    provider: GitProvider,
    selfHostedParams?: {
      gitlabUrl: string;
      clientId: string;
      clientSecret: string;
    }
  ) => void;
  handleDisconnect: () => Promise<void>;

  // Repository registration
  handleRegisterRepository: () => Promise<void>;

  // Loading states
  isLoadingConnections: boolean;
  isLoadingRepositories: boolean;
  isAuthenticating: boolean;
  isSubmitting: boolean;

  // Data
  connections: unknown[];
  repositories: GitRepository[];

  // Errors and messages
  error: string | null;
  message: { type: "success" | "error"; text: string } | null;
}

export function useRepositoryCreatePage(): RepositoryCreatePageState {
  const location = useLocation();
  const [selectedProvider, setSelectedProvider] =
    useState<GitProvider>("github");
  const [selectedRepository, setSelectedRepository] =
    useState<GitRepository | null>(null);

  const {
    isAuthenticating,
    isSubmitting,
    message,
    handleOAuthConnect,
    handleSubmitRepository,
  } = useRepositoryRegistration();

  const {
    connections,
    repositories,
    isLoadingConnections,
    isLoadingRepositories,
    error,
    isConnected: checkIsConnected,
    loadRepositories,
    disconnect,
    refreshConnections,
  } = useGitProvider();

  const isConnected = checkIsConnected(selectedProvider);

  // OAuthコールバックからの遷移を検出
  const locationState = location.state as {
    oauthSuccess?: boolean;
    provider?: string;
  } | null;

  // OAuth成功後にリポジトリ一覧を読み込む
  useEffect(() => {
    if (locationState?.oauthSuccess && locationState?.provider) {
      refreshConnections();
    }
  }, [locationState, refreshConnections]);

  // プロバイダーが接続済みの場合、リポジトリ一覧を読み込む
  useEffect(() => {
    if (isConnected) {
      loadRepositories(selectedProvider);
    }
  }, [selectedProvider, isConnected, loadRepositories, connections]);

  const handleConnect = useCallback(
    (
      provider: GitProvider,
      selfHostedParams?: {
        gitlabUrl: string;
        clientId: string;
        clientSecret: string;
      }
    ) => {
      if (selfHostedParams) {
        handleOAuthConnect(provider, selfHostedParams);
      } else {
        handleOAuthConnect(provider);
      }
    },
    [handleOAuthConnect]
  );

  const handleDisconnect = useCallback(async () => {
    await disconnect(selectedProvider);
    setSelectedRepository(null);
  }, [disconnect, selectedProvider]);

  const handleRegisterRepository = useCallback(async () => {
    if (!selectedRepository) return;

    await handleSubmitRepository({
      provider: selectedProvider,
      url: selectedRepository.cloneUrl,
    });
  }, [selectedRepository, selectedProvider, handleSubmitRepository]);

  return {
    selectedProvider,
    setSelectedProvider,
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
    connections,
    repositories,
    error,
    message,
  };
}
