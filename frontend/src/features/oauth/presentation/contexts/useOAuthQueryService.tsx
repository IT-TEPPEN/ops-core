import { useContext } from "react";
import { OAuthQueryService } from "../../application";
import { OAuthQueryServiceContext } from "./OAuthQueryServiceContext";

/**
 * Hook to access OAuthQueryService from context.
 */
export function useOAuthQueryService(): OAuthQueryService {
  const context = useContext(OAuthQueryServiceContext);

  if (!context) {
    throw new Error(
      "useOAuthQueryService must be used within OAuthQueryServiceProvider"
    );
  }

  return context;
}
