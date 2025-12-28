import { Route } from "react-router-dom";

import { DocumentGuard } from "./DocumentGuard";
import { DocumentDetailPage } from "./DocumentDetailPage";
import { DocumentListPage } from "./DocumentListPage";
import { DocumentVersionHistoryPage } from "./DocumentVersionHistoryPage";
import { DocumentViewPage } from "./DocumentViewPage";
import { ProtectedRoute } from "@/features/common/components/ProtectedRoute";
import ExecutionRecordPage from "../ExecutionRecordPage";

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
