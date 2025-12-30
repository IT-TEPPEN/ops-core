import { useContext } from "react";
import { OAuthCommandService } from "../../application";
import { OAuthCommandServiceContext } from "./OAuthCommandServiceContext";

/**
 * Hook to access OAuthCommandService from context.
 */
export function useOAuthCommandService(): OAuthCommandService {
  const context = useContext(OAuthCommandServiceContext);
  if (!context) {
    throw new Error(
      "useOAuthCommandService must be used within OAuthCommandServiceProvider"
    );
  }
  return context;
}
