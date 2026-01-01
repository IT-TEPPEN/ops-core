import { ProtectedRoute } from "@/features/common/components/ProtectedRoute";
import { Route } from "react-router-dom";
import { RepositoryListPage } from "./RepositoryListPage";
import { RepositoryDetailPage } from "./RepositoryDetailPage";
import { RepositoryIdGuard } from "./RepositoryGuard";

export function RepositoryRoutes() {
  return (
    <Route path="repositories">
      <Route
        index
        element={
          <ProtectedRoute>
            <RepositoryListPage />
          </ProtectedRoute>
        }
      />
      <Route path=":repoId">
        <Route
          index
          element={
            <ProtectedRoute>
              <RepositoryIdGuard element={RepositoryDetailPage} />
            </ProtectedRoute>
          }
        />
      </Route>
    </Route>
  );
}
