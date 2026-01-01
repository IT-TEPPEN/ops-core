// OAuth feature public API
// Following ADR 0021 - Pattern 1 (Simplified structure)

// Types
export type { GitProvider, GitRepository } from "./types";

// Services
export type {
  OAuthQueryService,
  OAuthCommandService,
} from "./application/services";

// Contexts
export {
  OAuthQueryServiceProvider,
  useOAuthQueryService,
  OAuthCommandServiceProvider,
  useOAuthCommandService,
  OAuthConnection,
  ConnectionList,
} from "./presentation";

// Hooks
export { useGitProvider } from "./hooks";
