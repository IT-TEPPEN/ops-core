import { TemporaryInfo, Token } from "../entity";

export interface SessionRepository {
  saveToken(token: Token): Promise<void>;
  getToken(): Promise<Token | null>;

  saveTemporaryInfo(info: TemporaryInfo): Promise<void>;
  getTemporaryInfo(): Promise<TemporaryInfo | null>;
  removeTemporaryInfo(): Promise<void>;
}
