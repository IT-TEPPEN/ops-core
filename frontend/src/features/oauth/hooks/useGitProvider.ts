import { useState, useEffect, useCallback } from "react";
import type { GitProvider, GitRepository } from "../types";
import {
  useOAuthQueryService,
  useOAuthCommandService,
} from "../presentation/contexts";
import { Connection } from "../application/dto";

interface UseGitProviderState {
  connections: Connection[];
  repositories: GitRepository[];
  isLoadingRepositories: boolean;
  error: string | null;
}

interface UseGitProviderReturn extends UseGitProviderState {
  isConnected: (provider: GitProvider) => boolean;
  getConnection: (provider: GitProvider) => Connection | undefined;
  loadRepositories: (provider: GitProvider) => Promise<void>;
  disconnect: (provider: GitProvider) => Promise<void>;
  refreshConnections: () => Promise<void>;
}

/**
 * Gitプロバイダーとの接続状態を管理するフック
 */
export function useGitProvider(): UseGitProviderReturn {
  const oauthQueryService = useOAuthQueryService();
  const oauthCommandService = useOAuthCommandService();

  const [state, setState] = useState<UseGitProviderState>({
    connections: [],
    repositories: [],
    isLoadingRepositories: false,
    error: null,
  });

  const refreshConnections = useCallback(async () => {
    setState((prev) => ({ ...prev, error: null }));
    try {
      const connections = await oauthQueryService.listConnections();
      setState((prev) => ({
        ...prev,
        connections,
      }));
    } catch (err) {
      setState((prev) => ({
        ...prev,
        error:
          err instanceof Error
            ? err.message
            : "Failed to load OAuth connections",
      }));
    }
  }, [oauthQueryService]);

  const isConnected = useCallback(
    (provider: GitProvider): boolean => {
      return state.connections.some((conn) => conn.provider === provider);
    },
    [state.connections]
  );

  const getConnection = useCallback(
    (provider: GitProvider): Connection | undefined => {
      return state.connections.find((conn) => conn.provider === provider);
    },
    [state.connections]
  );

  const loadRepositories = useCallback(
    async (provider: GitProvider) => {
      setState((prev) => ({
        ...prev,
        isLoadingRepositories: true,
        repositories: [],
        error: null,
      }));
      try {
        const connection = state.connections.find(
          (conn) => conn.provider === provider
        );
        if (!connection) {
          throw new Error(
            `No OAuth connection found for provider: ${provider}`
          );
        }
        const repositories = await oauthQueryService.listRepositories(
          connection.id
        );
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
    },
    [oauthQueryService, state.connections]
  );

  const disconnect = useCallback(
    async (provider: GitProvider) => {
      try {
        await oauthCommandService.disconnect(provider);
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
    [refreshConnections, oauthCommandService]
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
