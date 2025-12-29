import { useState, useEffect, useCallback } from "react";
import {
  gitProviderApi,
  GitProvider,
  OAuthConnection,
  GitRepository,
} from "@/shared/api/gitProviderApi";

interface UseGitProviderState {
  connections: OAuthConnection[];
  repositories: GitRepository[];
  isLoadingConnections: boolean;
  isLoadingRepositories: boolean;
  error: string | null;
}

interface UseGitProviderReturn extends UseGitProviderState {
  isConnected: (provider: GitProvider) => boolean;
  getConnection: (provider: GitProvider) => OAuthConnection | undefined;
  loadRepositories: (provider: GitProvider) => Promise<void>;
  disconnect: (provider: GitProvider) => Promise<void>;
  refreshConnections: () => Promise<void>;
}

/**
 * Gitプロバイダーとの接続状態を管理するフック
 */
export function useGitProvider(): UseGitProviderReturn {
  const [state, setState] = useState<UseGitProviderState>({
    connections: [],
    repositories: [],
    isLoadingConnections: true,
    isLoadingRepositories: false,
    error: null,
  });

  const refreshConnections = useCallback(async () => {
    setState((prev) => ({ ...prev, isLoadingConnections: true, error: null }));
    try {
      const connections = await gitProviderApi.listConnections();
      setState((prev) => ({
        ...prev,
        connections,
        isLoadingConnections: false,
      }));
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoadingConnections: false,
        error:
          err instanceof Error
            ? err.message
            : "Failed to load OAuth connections",
      }));
    }
  }, []);

  const isConnected = useCallback(
    (provider: GitProvider): boolean => {
      return state.connections.some((conn) => conn.provider === provider);
    },
    [state.connections]
  );

  const getConnection = useCallback(
    (provider: GitProvider): OAuthConnection | undefined => {
      return state.connections.find((conn) => conn.provider === provider);
    },
    [state.connections]
  );

  const loadRepositories = useCallback(async (provider: GitProvider) => {
    setState((prev) => ({
      ...prev,
      isLoadingRepositories: true,
      repositories: [],
      error: null,
    }));
    try {
      const repositories = await gitProviderApi.listRepositories(provider);
      setState((prev) => ({
        ...prev,
        repositories,
        isLoadingRepositories: false,
      }));
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoadingRepositories: false,
        error:
          err instanceof Error ? err.message : "Failed to load repositories",
      }));
    }
  }, []);

  const disconnect = useCallback(
    async (provider: GitProvider) => {
      try {
        await gitProviderApi.disconnect(provider);
        await refreshConnections();
        setState((prev) => ({ ...prev, repositories: [] }));
      } catch (err) {
        setState((prev) => ({
          ...prev,
          error:
            err instanceof Error
              ? err.message
              : "Failed to disconnect provider",
        }));
      }
    },
    [refreshConnections]
  );

  useEffect(() => {
    refreshConnections();
  }, [refreshConnections]);

  return {
    ...state,
    isConnected,
    getConnection,
    loadRepositories,
    disconnect,
    refreshConnections,
  };
}
