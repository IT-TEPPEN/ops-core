import { Routes, Route } from "react-router-dom";
import BlogPage from "../pages/BlogPage";
import RepositoryListPage from "../pages/RepositoryListPage";
import RepositoryCreatePage from "../pages/RepositoryCreatePage";
import RepositoryDetailPage from "../pages/RepositoryDetailPage";
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

function App() {
  return (
    <AuthProvider>
      <DiProviders>
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
          <Header />

          {/* Page Content Area */}
          <main className="max-w-5xl mx-auto p-4">
            <Routes>
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
                        <RepositoryDetailPage />
                      </ProtectedRoute>
                    }
                  />
                  <Route
                    path="files/:filePath"
                    element={
                      <ProtectedRoute>
                        <BlogPage />
                      </ProtectedRoute>
                    }
                  />
                </Route>
              </Route>
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
