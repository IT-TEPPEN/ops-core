import { describe, expect, it } from "vitest";
import { TokenImpl } from "./Token";

describe("Token Entity", () => {
  /** 
   * A sample JWT payload:
{
  "sub": "40bdbcd1-0285-49be-ba7f-c0d841553e77",
  "email": "yk.iteppen@gmail.com",
  "name": "上山泰史",
  "iss": "opscore",
  "exp": 1767366316,
  "iat": 1767365416
}
  */
  const sampleAccessToken =
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MGJkYmNkMS0wMjg1LTQ5YmUtYmE3Zi1jMGQ4NDE1NTNlNzciLCJlbWFpbCI6InlrLml0ZXBwZW5AZ21haWwuY29tIiwibmFtZSI6IuS4iuWxseazsOWPsiIsImlzcyI6Im9wc2NvcmUiLCJleHAiOjE3NjczNjYzMTYsImlhdCI6MTc2NzM2NTQxNn0.W8BZ4JAqG7kaevLLLL1anzMBMp4N3y-hyyUkBEjYpxM";
  const sampleRefreshToken = "sample_refresh_token";

  it("should create a Token from JWT", () => {
    const token = TokenImpl.fromJwt(sampleAccessToken, sampleRefreshToken);

    expect(token.accessToken).toBe(sampleAccessToken);
    expect(token.refreshToken).toBe(sampleRefreshToken);
    expect(token.userId).toBe("40bdbcd1-0285-49be-ba7f-c0d841553e77");
    expect(token.email).toBe("yk.iteppen@gmail.com");
    expect(token.name).toBe("上山泰史");
    expect(token.expiresAt.getTime()).toBe(1767366316 * 1000);
  });
});
