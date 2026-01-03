import z from "zod";

export interface Token {
  get accessToken(): string;
  get refreshToken(): string;
  get userId(): string;
  get email(): string;
  get name(): string;
  get expiresAt(): Date;
}

const jwtPayloadSchema = z.object({
  sub: z.string(),
  email: z.email(),
  name: z.string(),
  exp: z.number(),
});

type JwtPayloadSchema = z.infer<typeof jwtPayloadSchema>;

export class TokenImpl implements Token {
  private constructor(
    private readonly _accessToken: string,
    private readonly _refreshToken: string,
    private readonly _userId: string,
    private readonly _email: string,
    private readonly _name: string,
    private readonly _expiresAt: Date
  ) {}

  // Getters (implementing Token interface)
  get accessToken(): string {
    return this._accessToken;
  }
  get refreshToken(): string {
    return this._refreshToken;
  }
  get userId(): string {
    return this._userId;
  }
  get email(): string {
    return this._email;
  }
  get name(): string {
    return this._name;
  }
  get expiresAt(): Date {
    return this._expiresAt;
  }

  // Static factory method for JWT conversion
  static fromJwt(accessToken: string, refreshToken: string): TokenImpl {
    const payload = TokenImpl.decodeJwtPayload(accessToken);

    return new TokenImpl(
      accessToken,
      refreshToken,
      payload.sub,
      payload.email,
      payload.name,
      new Date(payload.exp * 1000)
    );
  }

  // Private JWT decoder
  static decodeJwtPayload(jwt: string): JwtPayloadSchema {
    try {
      const parts = jwt.split(".");
      if (parts.length !== 3) {
        throw new Error("Invalid JWT format");
      }

      const payload = parts[1];
      const decoded = new TextDecoder("utf-8").decode(
        Uint8Array.from(atob(payload), (c) => c.charCodeAt(0))
      );
      return jwtPayloadSchema.parse(JSON.parse(decoded));
    } catch (error) {
      throw new Error(
        `Failed to decode JWT: ${
          error instanceof Error ? error.message : "unknown error"
        }`
      );
    }
  }
}
