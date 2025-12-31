/**
 * OAuth feature type definitions.
 * Following Pattern 1 - Simplified structure.
 */

// Re-export types from shared API for convenience
// These types are defined in @/shared/api/gitProviderApi
export type { GitProvider, GitRepository } from "@/shared/api/gitProviderApi";

/**
 * Git provider state for UI.
 */
export interface UseGitProviderState {
  connections: any[];
  repositories: any[];
  isLoadingConnections: boolean;
  isLoadingRepositories: boolean;
  error: string | null;
}

/**
 * Git provider hook return type.
 */
export interface UseGitProviderReturn extends UseGitProviderState {
  isConnected: (provider: any) => boolean;
  getConnection: (provider: any) => any | undefined;
  loadRepositories: (provider: any) => Promise<void>;
  disconnect: (provider: any) => Promise<void>;
  refreshConnections: () => Promise<void>;
}
