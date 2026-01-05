import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useSignOutUsecase } from "../contexts";

export function useSignOut() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const signOutUsecase = useSignOutUsecase();
  const navigate = useNavigate();

  const signOut = async () => {
    setIsLoading(true);
    setError(null);

    try {
      await signOutUsecase.execute();
      // Redirect to home page after sign out
      navigate("/");
      // Reload to clear any cached state
      window.location.reload();
    } catch (err) {
      setError(err instanceof Error ? err : new Error("Failed to sign out"));
    } finally {
      setIsLoading(false);
    }
  };

  return {
    signOut,
    isLoading,
    error,
  };
}
