import { ProtectedRoute } from "@/features/common/components/ProtectedRoute";
import { Route } from "react-router-dom";
import { TypedParamsGuard } from "@/shared/components";
import { RepositoryListPage } from "./RepositoryListPage";
import { DocumentPreviewPage } from "./DocumentPreviewPage";
import { RepositoryCreatePage } from "./RepositoryCreatePage";
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
      <Route
        path="new"
        element={
          <ProtectedRoute>
            <RepositoryCreatePage />
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
        <Route
          path="files/:filePath"
          element={
            <ProtectedRoute>
              <TypedParamsGuard
                element={DocumentPreviewPage}
                pathParams={[
                  {
                    name: "repoId",
                    required: true,
                  },
                  {
                    name: "filePath",
                    required: true,
                  },
                ]}
                queryParams={undefined}
                fallbackPath={"/repositories"}
              />
            </ProtectedRoute>
          }
        />
      </Route>
    </Route>
  );
}
