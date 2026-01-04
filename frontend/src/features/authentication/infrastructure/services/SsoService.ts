import { V1ApiClient } from "@/shared/api";
import { AuthenticationService } from "../../application/services/AuthenticationService";
import {
  GetProviderLoginUrlDto,
  IdentitiesDto,
  IdentityDto,
  LoginUriDto,
  TokenDto,
  UnlinkIdentityDto,
  ValidateCodeAndGetTokenDto,
} from "../../application/dto";
import { AxiosError } from "axios";

interface LoginUrlResponse {
  auth_url: string;
  state: string;
}

interface AuthCallbackRequest {
  code: string;
  remember_me: boolean;
  state: string;
}

interface TokenResponse {
  token: string;
  refresh_token: string;
}

interface IdentityResponse {
  id: string;
  provider: string;
  email: string;
  name: string;
  picture: string;
  linked_at: string;
  last_used_at: string;
}

interface IdentitiesResponse {
  identities: IdentityResponse[];
}

export class SsoService extends V1ApiClient implements AuthenticationService {
  constructor() {
    super("/auth");
  }

  async getLoginUrl(dto: GetProviderLoginUrlDto): Promise<LoginUriDto> {
    const response = await this.get<LoginUrlResponse>(
      `/${dto.providerName}/login`
    );

    return {
      authUrl: response.auth_url,
      state: response.state,
    };
  }

  async validateCodeAndGetToken(
    dto: ValidateCodeAndGetTokenDto
  ): Promise<TokenDto> {
    const response = await this.post<TokenResponse, AuthCallbackRequest>(
      `/${dto.provider}/callback`,
      { code: dto.code, state: dto.state, remember_me: dto.rememberMe }
    );

    return {
      accessToken: response.token,
      refreshToken: response.refresh_token,
    };
  }

  async getIdentity(): Promise<IdentityDto | null> {
    try {
      const response = await this.get<IdentityResponse>("/identity");

      return {
        id: response.id,
        provider: response.provider,
        email: response.email,
        name: response.name,
        pictureUrl: response.picture,
        linkedAt: response.linked_at,
        lastUsedAt: response.last_used_at,
      };
    } catch (error) {
      if (error instanceof AxiosError && error.response?.status === 401) {
        return null;
      }
      throw error;
    }
  }

  async listIdentities(): Promise<IdentitiesDto> {
    const response = await this.get<IdentitiesResponse>("/identities");

    return {
      identities: response.identities.map((identity) => ({
        id: identity.id,
        provider: identity.provider,
        email: identity.email,
        name: identity.name,
        pictureUrl: identity.picture,
        linkedAt: identity.linked_at,
        lastUsedAt: identity.last_used_at,
      })),
    };
  }

  async unlinkIdentity(dto: UnlinkIdentityDto): Promise<void> {
    await this.delete<void>(`/identities/${dto.identityId}`);
  }

  async logout(): Promise<void> {
    await this.post<void, Record<string, never>>("/logout", {});
  }
}
