import {
  TemporaryInfo,
  TemporaryInfoImpl,
  Token,
  TokenImpl,
} from "../../domain/entity";
import { SessionRepository } from "../../domain/repository";

export class SessionLocalStorageRepository implements SessionRepository {
  private readonly ACCESS_TOKEN_KEY = "access_token";
  private readonly REFRESH_TOKEN_KEY = "refresh_token";
  private readonly TEMPORARY_INFO_KEY = "temporary_info";

  async saveToken(session: Token): Promise<void> {
    localStorage.setItem(this.ACCESS_TOKEN_KEY, session.accessToken);
    localStorage.setItem(this.REFRESH_TOKEN_KEY, session.refreshToken);
  }

  async getToken(): Promise<Token | null> {
    const accessToken = localStorage.getItem(this.ACCESS_TOKEN_KEY);
    const refreshToken = localStorage.getItem(this.REFRESH_TOKEN_KEY);

    if (!accessToken || !refreshToken) {
      return null;
    }

    return TokenImpl.fromJwt(accessToken, refreshToken);
  }

  async removeToken(): Promise<void> {
    localStorage.removeItem(this.ACCESS_TOKEN_KEY);
    localStorage.removeItem(this.REFRESH_TOKEN_KEY);
  }

  async saveTemporaryInfo(info: TemporaryInfo): Promise<void> {
    localStorage.setItem(this.TEMPORARY_INFO_KEY, info.toJSONString());
  }

  async getTemporaryInfo(): Promise<TemporaryInfo | null> {
    const jsonString = localStorage.getItem(this.TEMPORARY_INFO_KEY);
    if (!jsonString) {
      return null;
    }
    return TemporaryInfoImpl.fromJSONString(jsonString);
  }

  async removeTemporaryInfo(): Promise<void> {
    localStorage.removeItem(this.TEMPORARY_INFO_KEY);
  }
}
