import { useEffect } from "react";
import { useNavigate, useSearchParams, useParams } from "react-router-dom";
import {
  useAuthenticationService,
  useSessionRepository,
} from "@/features/authentication/infrastructure/contexts";
import { TokenImpl } from "@/features/authentication/domain/entity";

export default function AuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const authenticationService = useAuthenticationService();
  const sessionRepository = useSessionRepository();

  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const { provider } = useParams<{ provider: string }>();

  useEffect(() => {
    const handleCallback = async () => {
      if (!code || !state || !provider) {
        console.error("Missing code, state, or provider parameter");
        navigate("/login", { replace: true });
        return;
      }

      const temporarySession = await sessionRepository.getTemporaryInfo();
      const rememberMe = sessionStorage.getItem("auth_remember_me") === "true";

      if (
        state !== temporarySession?.state ||
        provider !== temporarySession?.provider
      ) {
        console.error("State or provider mismatch - potential CSRF attack");
        navigate("/login", { replace: true });
        return;
      }

      try {
        const tokensDto = await authenticationService.validateCodeAndGetToken({
          provider,
          code,
          state,
          rememberMe,
        });

        const tokens = TokenImpl.fromJwt(
          tokensDto.accessToken,
          tokensDto.refreshToken
        );

        await sessionRepository.saveToken(tokens);
        await sessionRepository.removeTemporaryInfo();

        navigate(temporarySession.from || "/", { replace: true });
      } catch (error) {
        console.error("Authentication error:", error);
        navigate("/login", { replace: true });
      }
    };

    handleCallback();
  }, [
    navigate,
    provider,
    code,
    state,
    sessionRepository,
    authenticationService,
  ]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-4 text-gray-600">Completing sign in...</p>
      </div>
    </div>
  );
}
