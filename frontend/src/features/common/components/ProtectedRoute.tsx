import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../../../app/hooks/useAuth";
import { LoadingSpinner } from "@/ui";
import { useGetUserIdentityUsecase } from "@/features/authentication/presentation/contexts";
import { useQuery } from "@tanstack/react-query";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const getUserIdentityUsecase = useGetUserIdentityUsecase();
  const query = useQuery({
    queryKey: ["UserIdentity"],
    queryFn: async () => getUserIdentityUsecase.execute(),
  });
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  if (query.isLoading) {
    return <LoadingSpinner message="Loading..." />;
  }

  if (query.isError) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          Error loading user identity: {`${query.error}`}
        </div>
      </div>
    );
  }

  if (!query.data) {
    // Redirect to login page, saving the attempted location
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}
