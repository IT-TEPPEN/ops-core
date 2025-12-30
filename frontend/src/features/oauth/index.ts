// OAuth feature public API
// Following ADR 0021 - Pattern 1 (Simplified structure)

// Types
export type {
  GitProvider,
  OAuthConnection as OAuthConnectionType,
  GitRepository,
  UseGitProviderState,
  UseGitProviderReturn,
} from "./types";

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
} from "./presentation/contexts";

// Hooks
export { useGitProvider } from "./hooks";

// Components
export { OAuthConnection, ConnectionStatus } from "./components";
