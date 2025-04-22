import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

interface Repository {
  id: string;
  name: string;
  url: string;
  createdAt: string;
  updatedAt: string;
}

function RepositoriesPage() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newRepoUrl, setNewRepoUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitMessage, setSubmitMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api";

  // Fetch repositories on component mount
  useEffect(() => {
    fetchRepositories();
  }, []);

  const fetchRepositories = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${apiUrl}/repositories`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setRepositories(data.repositories);
    } catch (err) {
      setError("Failed to load repositories. Please try again later.");
      console.error("Error fetching repositories:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitMessage(null);

    try {
      const response = await fetch(`${apiUrl}/repositories`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: newRepoUrl }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to register repository");
      }

      setSubmitMessage({
        type: "success",
        text: "Repository registered successfully!",
      });
      setNewRepoUrl("");
      fetchRepositories(); // Refresh the list
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      setSubmitMessage({ type: "error", text: message });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Repository Management</h1>

      {/* Registration Form */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Register New Repository</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="repoUrl" className="block text-sm font-medium mb-1">
              Repository URL
            </label>
            <input
              id="repoUrl"
              type="text"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              placeholder="https://github.com/username/repo.git"
              value={newRepoUrl}
              onChange={(e) => setNewRepoUrl(e.target.value)}
              required
            />
          </div>
          <button
            type="submit"
            disabled={isSubmitting}
            className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
              isSubmitting ? "opacity-50 cursor-not-allowed" : ""
            }`}
          >
            {isSubmitting ? "Registering..." : "Register Repository"}
          </button>
        </form>

        {submitMessage && (
          <div
            className={`mt-4 p-3 rounded ${
              submitMessage.type === "success"
                ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
            }`}
          >
            {submitMessage.text}
          </div>
        )}
      </div>

      {/* Repository List */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Registered Repositories</h2>

        {isLoading && <p className="text-gray-500">Loading repositories...</p>}

        {error && (
          <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
            {error}
          </div>
        )}

        {!isLoading && !error && repositories.length === 0 && (
          <p className="text-gray-500">No repositories registered yet.</p>
        )}

        {repositories.length > 0 && (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    URL
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Created
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                {repositories.map((repo) => (
                  <tr key={repo.id}>
                    <td className="px-6 py-4 whitespace-nowrap">{repo.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {repo.url}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {new Date(repo.createdAt).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <Link
                        to={`/repositories/${repo.id}`}
                        className="text-blue-500 hover:text-blue-700 font-medium"
                      >
                        Manage Files
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default RepositoriesPage;
