import type { GitProvider } from "../../types";

/**
 * OAuth Command Service interface for write operations (DELETE).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface OAuthCommandService {
  /**
   * Disconnect from a provider.
   */
  disconnect(provider: GitProvider): Promise<void>;
}
