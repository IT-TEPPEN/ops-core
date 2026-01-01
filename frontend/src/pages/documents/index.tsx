import { Route } from "react-router-dom";

import { DocumentGuard } from "./DocumentGuard";
import { DocumentDetailPage } from "./DocumentDetailPage";
import { DocumentListPage } from "./DocumentListPage";
import { DocumentPublishPage } from "./DocumentPublishPage";
import { DocumentVersionHistoryPage } from "./DocumentVersionHistoryPage";
import { DocumentViewPage } from "./DocumentViewPage";
import { ProtectedRoute } from "@/features/common/components/ProtectedRoute";
import ExecutionRecordPage from "../ExecutionRecordPage";
import { TypedParamsGuard } from "@/shared/components";
import { DocumentPreviewPage } from "./DocumentPreviewPage";

export function DocumentRoutes() {
  return (
    <Route path="documents">
      <Route
        index
        element={
          <ProtectedRoute>
            <DocumentListPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="publish"
        element={
          <ProtectedRoute>
            <DocumentPublishPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="preview/"
        element={
          <ProtectedRoute>
            <TypedParamsGuard
              element={DocumentPreviewPage}
              pathParams={undefined}
              queryParams={[
                {
                  name: "connection_id",
                  required: true,
                },
                {
                  name: "repository_full_name",
                  required: true,
                },
                {
                  name: "path",
                  required: true,
                },
              ]}
              fallbackPath={"/repositories"}
            />
          </ProtectedRoute>
        }
      />
      <Route path=":docId">
        <Route
          index
          element={
            <ProtectedRoute>
              <DocumentGuard element={DocumentDetailPage} />
            </ProtectedRoute>
          }
        />
        <Route
          path="view"
          element={
            <ProtectedRoute>
              <DocumentGuard element={DocumentViewPage} />
            </ProtectedRoute>
          }
        />
        <Route path="execute">
          <Route
            path=":recordId"
            element={
              <ProtectedRoute>
                <ExecutionRecordPage />
              </ProtectedRoute>
            }
          />
        </Route>
        <Route
          path="versions"
          element={
            <ProtectedRoute>
              <DocumentGuard element={DocumentVersionHistoryPage} />
            </ProtectedRoute>
          }
        />
      </Route>
    </Route>
  );
}
