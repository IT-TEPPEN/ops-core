import { useState, useEffect } from "react";
import { Routes, Route, Link } from "react-router-dom";
import BlogPage from "./pages/BlogPage";
import RepositoriesPage from "./pages/RepositoriesPage";
import RepositoryDetailPage from "./pages/RepositoryDetailPage";
import DocumentListPage from "./pages/DocumentListPage";
import DocumentDetailPage from "./pages/DocumentDetailPage";
import DocumentViewPage from "./pages/DocumentViewPage";
import DocumentVersionHistoryPage from "./pages/DocumentVersionHistoryPage";

// Define a simple Home component for the root path
function HomePage() {
  const [message, setMessage] = useState("Loading...");

  useEffect(() => {
    const apiHost = import.meta.env.VITE_API_HOST;
    const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api";

    console.log(`Fetching data from: ${apiUrl}`);

    fetch(apiUrl)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => setMessage(data.message))
      .catch((error) => {
        console.error("Error fetching data:", error);
        setMessage("Failed to load message from backend.");
      });
  }, []);

  return (
    <div className="text-center p-8">
      <h1 className="text-3xl font-bold mb-4">OpsCore Documentation System</h1>
      <p className="mb-8 text-lg">
        A system for managing operational procedure documents from external
        repositories
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto">
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow">
          <h2 className="text-xl font-bold mb-2">Repository Management</h2>
          <p className="mb-4 text-gray-600 dark:text-gray-300">
            Register external Git repositories and select markdown files to
            display as documentation.
          </p>
          <Link
            to="/repositories"
            className="inline-block px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
          >
            Manage Repositories
          </Link>
        </div>

        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow">
          <h2 className="text-xl font-bold mb-2">Documentation Viewer</h2>
          <p className="mb-4 text-gray-600 dark:text-gray-300">
            View the selected markdown files as formatted documentation pages.
          </p>
          <Link
            to="/blog"
            className="inline-block px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition"
          >
            View Documentation
          </Link>
        </div>
      </div>

      <div className="mt-12 text-gray-500">
        Backend status: <span className="font-semibold">{message}</span>
      </div>
    </div>
  );
}

function App() {
  return (
    // Apply base layout and background colors using Tailwind
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Navigation Bar */}
      <nav className="bg-white dark:bg-gray-800 shadow-md sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-4">
          <div className="flex justify-center items-center h-16">
            <ul className="flex space-x-6">
              <li>
                <Link
                  to="/"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Home
                </Link>
              </li>
              <li>
                <Link
                  to="/repositories"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Repositories
                </Link>
              </li>
              <li>
                <Link
                  to="/documents"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Documents
                </Link>
              </li>
              <li>
                <Link
                  to="/blog"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Documentation
                </Link>
              </li>
            </ul>
          </div>
        </div>
      </nav>

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
          <Route path="/documents/:docId/view" element={<DocumentViewPage />} />
          <Route path="/documents/:docId/versions" element={<DocumentVersionHistoryPage />} />
          <Route path="/blog" element={<BlogPage />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
