import { Routes, Route } from "react-router-dom";
import AuthCallbackPage from "../pages/AuthCallbackPage";
import { LoginPage } from "../pages/LoginPage";
import AccountSettingsPage from "../pages/AccountSettingsPage";
import { DocumentRoutes } from "@/pages/documents";
import { Header } from "../features/common/components/Layout/Header";
import { HomePage } from "../pages/Home";
import { DiProviders } from "./providers";
import { NotificationCard } from "../features/notification";
import { ProtectedRoute } from "../features/common/components/ProtectedRoute";
import { MarkdownTestEditor } from "@/features/markdown";
import { OAuthRoutes } from "@/pages/oauth";
import { GroupRoutes } from "@/pages/groups";

function App() {
  return (
    <DiProviders>
      <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <div className="h-[4rem]">
          <Header />
        </div>

        {/* Page Content Area */}
        <main className="h-[calc(100vh-4rem)] overflow-auto">
          <Routes>
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <HomePage />
                </ProtectedRoute>
              }
            />
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

            {OAuthRoutes()}
            {DocumentRoutes()}
            {GroupRoutes()}
          </Routes>
        </main>

        <div className="fixed bottom-4 right-4">
          <NotificationCard />
        </div>
      </div>
    </DiProviders>
  );
}

export default App;
