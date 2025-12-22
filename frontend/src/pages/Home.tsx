import { Link } from "react-router-dom";

export function HomePage() {
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
    </div>
  );
}
