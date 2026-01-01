import { ProtectedRoute } from "@/features/common/components";
import { Route } from "react-router-dom";
import GroupListPage from "./GroupListPage";
import GroupDetailPage from "./GroupDetailPage";

export function GroupRoutes() {
  return (
    <Route path="/groups">
      <Route
        index
        element={
          <ProtectedRoute>
            <GroupListPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/groups/:groupId"
        element={
          <ProtectedRoute>
            <GroupDetailPage />
          </ProtectedRoute>
        }
      />
    </Route>
  );
}
