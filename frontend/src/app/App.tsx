import { Routes, Route } from "react-router-dom";
import OAuthCallbackPage from "../pages/OAuthCallbackPage";
import AuthCallbackPage from "../pages/AuthCallbackPage";
import LoginPage from "../pages/LoginPage";
import AccountSettingsPage from "../pages/AccountSettingsPage";
import { DocumentRoutes } from "@/pages/documents";
import GroupListPage from "../pages/GroupListPage";
import GroupDetailPage from "../pages/GroupDetailPage";
import { Header } from "../features/common/components/Layout/Header";
import { HomePage } from "../pages/Home";
import { DiProviders } from "./providers";
import { NotificationCard } from "../features/notification";
import { AuthProvider } from "./contexts/AuthContext";
import { ProtectedRoute } from "../features/common/components/ProtectedRoute";
import { RepositoryRoutes } from "@/pages/repositories";
import { MarkdownTestEditor } from "@/features/markdown";

function App() {
  return (
    <AuthProvider>
      <DiProviders>
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
          <div className="h-16">
            <Header />
          </div>

          {/* Page Content Area */}
          <main className="h-[calc(100vh-4rem)] p-4">
            <Routes>
              <Route path="/markdown-test" element={<MarkdownTestEditor />} />
              <Route path="/login" element={<LoginPage />} />
              <Route
                path="/auth/:provider/callback"
                element={<AuthCallbackPage />}
              />
              <Route path="/oauth/callback" element={<OAuthCallbackPage />} />
              <Route
                path="/settings/account"
                element={
                  <ProtectedRoute>
                    <AccountSettingsPage />
                  </ProtectedRoute>
                }
              />

              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <HomePage />
                  </ProtectedRoute>
                }
              />
              {RepositoryRoutes()}
              {DocumentRoutes()}
              <Route
                path="/groups"
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
            </Routes>
          </main>

          <div className="fixed bottom-4 right-4">
            <NotificationCard />
          </div>
        </div>
      </DiProviders>
    </AuthProvider>
  );
}

export default App;
