import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";

interface Repository {
  id: string;
  name: string;
  url: string;
  createdAt: string;
  updatedAt: string;
}

interface FileNode {
  path: string;
  type: "file" | "dir";
}

function RepositoryDetailPage() {
  const { repoId } = useParams<{ repoId: string }>();

  const [repository, setRepository] = useState<Repository | null>(null);
  const [files, setFiles] = useState<FileNode[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);

  // Access token state
  const [accessToken, setAccessToken] = useState<string>("");
  const [isUpdatingToken, setIsUpdatingToken] = useState(false);
  const [tokenMessage, setTokenMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);
  const [needsToken, setNeedsToken] = useState(false);

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api";

  // Fetch repository details and files on component mount
  useEffect(() => {
    if (repoId) {
      Promise.all([fetchRepository(), fetchFiles()]);
    }
  }, [repoId]);

  const fetchRepository = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${apiUrl}/repositories/${repoId}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setRepository(data);
    } catch (err) {
      setError("Failed to load repository details. Please try again later.");
      console.error("Error fetching repository:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchFiles = async () => {
    setIsLoading(true);
    setFileError(null);
    setNeedsToken(false);

    try {
      const response = await fetch(`${apiUrl}/repositories/${repoId}/files`);

      if (response.status === 400) {
        // Check if this is an access token error
        const errorData = await response.json();
        if (errorData.code === "ACCESS_TOKEN_REQUIRED") {
          setNeedsToken(true);
          setFileError("Access token is required to list repository files");
          setFiles([]);
          setIsLoading(false);
          return;
        }
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setFiles(data.files);
    } catch (err) {
      setFileError("Failed to load repository files. Please try again later.");
      console.error("Error fetching files:", err);
    } finally {
      setIsLoading(false);
    }
  };

  // Handle access token update
  const handleTokenSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!accessToken.trim()) {
      setTokenMessage({
        type: "error",
        text: "Please enter an access token",
      });
      return;
    }

    setIsUpdatingToken(true);
    setTokenMessage(null);

    try {
      const response = await fetch(`${apiUrl}/repositories/${repoId}/token`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ accessToken }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to update access token");
      }

      setTokenMessage({
        type: "success",
        text: "Access token updated successfully!",
      });

      // Clear the form and refetch files
      setAccessToken("");
      setNeedsToken(false);
      fetchFiles();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      setTokenMessage({ type: "error", text: message });
    } finally {
      setIsUpdatingToken(false);
    }
  };

  // Filter files to only show markdown files
  const markdownFiles = files.filter(
    (file) => file.type === "file" && file.path.toLowerCase().endsWith(".md")
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Repository Details</h1>
        <Link
          to="/repositories"
          className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600"
        >
          Back to Repositories
        </Link>
      </div>

      {isLoading && !repository && (
        <p className="text-gray-500">Loading repository information...</p>
      )}

      {error && (
        <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
      )}

      {repository && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-2">{repository.name}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            <span className="font-medium">URL:</span> {repository.url}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            <span className="font-medium">Registered on:</span>{" "}
            {new Date(repository.createdAt).toLocaleString()}
          </p>
        </div>
      )}

      {/* Access Token Form */}
      {repository && (needsToken || fileError) && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-4">
            Repository Access Token
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
            This repository requires an access token to view files. Please enter
            a valid access token below.
          </p>

          <form onSubmit={handleTokenSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="accessToken"
                className="block text-sm font-medium mb-1"
              >
                Access Token
              </label>
              <input
                id="accessToken"
                type="password"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                placeholder="Enter GitHub personal access token"
                value={accessToken}
                onChange={(e) => setAccessToken(e.target.value)}
                required
              />
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                For GitHub repositories, create a personal access token with
                'repo' scope.
              </p>
            </div>

            <div>
              <button
                type="submit"
                disabled={isUpdatingToken}
                className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
                  isUpdatingToken ? "opacity-50 cursor-not-allowed" : ""
                }`}
              >
                {isUpdatingToken ? "Updating..." : "Update Access Token"}
              </button>
            </div>

            {tokenMessage && (
              <div
                className={`mt-4 p-3 rounded ${
                  tokenMessage.type === "success"
                    ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                    : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
                }`}
              >
                {tokenMessage.text}
              </div>
            )}
          </form>
        </div>
      )}

      {/* File Selection */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Select Markdown Files</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
          Select markdown files from the repository to display as documentation
          pages in OpsCore.
        </p>

        {fileError && !isLoading && (
          <div className="p-3 mb-4 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
            {fileError}
          </div>
        )}

        {isLoading && repository && (
          <p className="text-gray-500">Loading repository files...</p>
        )}

        {!isLoading && !fileError && markdownFiles.length === 0 ? (
          <div className="text-gray-500 dark:text-gray-400 mb-4">
            No markdown files found in this repository.
          </div>
        ) : (
          !needsToken &&
          !fileError && (
            <div className="overflow-y-auto max-h-96 border border-gray-200 dark:border-gray-700 rounded p-2">
              <table className="min-w-full">
                <thead>
                  <tr>
                    <th className="px-4 py-2 text-left text-sm font-medium text-gray-500 dark:text-gray-400">
                      File Path
                    </th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {markdownFiles.map((file) => (
                    <tr
                      key={file.path}
                      className="hover:bg-gray-100 dark:hover:bg-gray-700"
                    >
                      <td className="px-4 py-2 text-sm">{file.path}</td>
                      <td className="px-4 py-2 text-sm">
                        <Link
                          to={`/repositories/${
                            repository?.id
                          }/files/${encodeURIComponent(file.path)}`}
                          className="text-blue-500 hover:text-blue-700 font-medium"
                        >
                          Preview
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )
        )}
      </div>
    </div>
  );
}

export default RepositoryDetailPage;
