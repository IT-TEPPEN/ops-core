import type { OAuthCommandService } from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import {
  OAuthConnectRequest,
  OAuthConnectionResponse,
} from "../../application/services/OAuthCommandService";

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

  initiateOAuthFlow(_provider: string, _redirectUri: string): Promise<string> {
    throw new Error("Method not implemented.");
  }

  connectOAuth(request: OAuthConnectRequest): Promise<OAuthConnectionResponse> {
    return this.post<OAuthConnectionResponse, OAuthConnectRequest>(
      "/oauth/connect",
      request
    );
  }
}
