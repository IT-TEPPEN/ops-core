import type { OAuthCommandService } from "../../application";
import type { GitProvider } from "../../types";
import { V1ApiClient } from "@/shared/api/client";

/**
 * HTTP implementation of OAuthCommandService.
 * Handles all write operations (DELETE) for OAuth management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpOAuthCommandService
  extends V1ApiClient
  implements OAuthCommandService
{
  constructor() {
    super("");
  }

  async disconnect(provider: GitProvider): Promise<void> {
    await this.delete(`/auth/oauth/connections/${provider}`);
  }
}
