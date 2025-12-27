import { Routes, Route } from "react-router-dom";
import BlogPage from "./pages/BlogPage";
import RepositoriesPage from "./pages/RepositoriesPage";
import RepositoryDetailPage from "./pages/RepositoryDetailPage";
import DocumentListPage from "./pages/DocumentListPage";
import DocumentDetailPage from "./pages/DocumentDetailPage";
import DocumentViewPage from "./pages/DocumentViewPage";
import DocumentVersionHistoryPage from "./pages/DocumentVersionHistoryPage";
import ExecutionRecordPage from "./pages/ExecutionRecordPage";
import GroupListPage from "./pages/GroupListPage";
import GroupDetailPage from "./pages/GroupDetailPage";
import { Header } from "./components/Layout/Header";
import { HomePage } from "./pages/Home";
import { DiProviders } from "./providers";
import { NotificationCard } from "./features/notification/components";

function App() {
  return (
    <DiProviders>
      <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <Header />

        {/* Page Content Area */}
        <main className="max-w-5xl mx-auto p-4">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/repositories" element={<RepositoriesPage />} />
            <Route
              path="/repositories/:repoId"
              element={<RepositoryDetailPage />}
            />
            <Route path="/documents" element={<DocumentListPage />} />
            <Route path="/documents/:docId" element={<DocumentDetailPage />} />
            <Route
              path="/documents/:docId/view"
              element={<DocumentViewPage />}
            />
            <Route
              path="/documents/:docId/execute"
              element={<ExecutionRecordPage />}
            />
            <Route
              path="/documents/:docId/execute/:recordId"
              element={<ExecutionRecordPage />}
            />
            <Route
              path="/documents/:docId/versions"
              element={<DocumentVersionHistoryPage />}
            />
            <Route path="/groups" element={<GroupListPage />} />
            <Route path="/groups/:groupId" element={<GroupDetailPage />} />
            <Route path="/blog" element={<BlogPage />} />
          </Routes>
        </main>

        <div className="fixed bottom-4 right-4">
          <NotificationCard
            title="Success!!!"
            message="Your operation was completed successfully."
          />
        </div>
      </div>
    </DiProviders>
  );
}

export default App;
