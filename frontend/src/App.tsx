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
            <Route path="repositories">
              <Route index element={<RepositoriesPage />} />
              <Route path=":repoId">
                <Route index element={<RepositoryDetailPage />} />
                <Route path="files/:filePath" element={<BlogPage />} />
              </Route>
            </Route>
            <Route path="documents">
              <Route index element={<DocumentListPage />} />
              <Route path=":docId">
                <Route index element={<DocumentDetailPage />} />
                <Route path="view" element={<DocumentViewPage />} />
                <Route path="execute">
                  <Route path=":recordId" element={<ExecutionRecordPage />} />
                </Route>
                <Route
                  path="versions"
                  element={<DocumentVersionHistoryPage />}
                ></Route>
              </Route>
            </Route>
            <Route path="/groups" element={<GroupListPage />} />
            <Route path="/groups/:groupId" element={<GroupDetailPage />} />
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
