import { Routes, Route } from "react-router-dom";
import AuthCallbackPage from "../pages/AuthCallbackPage";
import LoginPage from "../pages/LoginPage";
import AccountSettingsPage from "../pages/AccountSettingsPage";
import { DocumentRoutes } from "@/pages/documents";
import { Header } from "../features/common/components/Layout/Header";
import { HomePage } from "../pages/Home";
import { DiProviders } from "./providers";
import { NotificationCard } from "../features/notification";
import { AuthProvider } from "./contexts/AuthContext";
import { ProtectedRoute } from "../features/common/components/ProtectedRoute";
import { RepositoryRoutes } from "@/pages/repositories";
import { MarkdownTestEditor } from "@/features/markdown";
import { OAuthRoutes } from "@/pages/oauth";
import { GroupRoutes } from "@/pages/groups";

function App() {
  return (
    <AuthProvider>
      <DiProviders>
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
          <div className="h-16">
            <Header />
          </div>

          {/* Page Content Area */}
          <main className="max-h-[calc(100vh-4rem)] overflow-auto p-4">
            <Routes>
              <Route path="/markdown-test" element={<MarkdownTestEditor />} />
              <Route path="/login" element={<LoginPage />} />
              <Route
                path="/auth/:provider/callback"
                element={<AuthCallbackPage />}
              />
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
              {OAuthRoutes()}
              {RepositoryRoutes()}
              {DocumentRoutes()}
              {GroupRoutes()}
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
