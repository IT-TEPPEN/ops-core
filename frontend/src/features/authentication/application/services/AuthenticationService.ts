import {
  GetProviderLoginUrlDto,
  IdentitiesDto,
  IdentityDto,
  LoginUriDto,
  TokenDto,
  UnlinkIdentityDto,
  ValidateCodeAndGetTokenDto,
} from "../dto";

export interface AuthenticationService {
  getLoginUrl(dto: GetProviderLoginUrlDto): Promise<LoginUriDto>;

  validateCodeAndGetToken(dto: ValidateCodeAndGetTokenDto): Promise<TokenDto>;

  getIdentity(): Promise<IdentityDto | null>;

  listIdentities(): Promise<IdentitiesDto>;

  unlinkIdentity(dto: UnlinkIdentityDto): Promise<void>;

  logout(): Promise<void>;
}
