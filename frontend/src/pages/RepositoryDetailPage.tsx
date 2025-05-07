import { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import {
  fetchRepositoryDetails,
  fetchRepositoryFiles,
  updateAccessToken,
  selectFiles,
} from "../api/repositories";
import ErrorMessage from "../ui/ErrorMessage";

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
  const navigate = useNavigate();

  const [repository, setRepository] = useState<Repository | null>(null);
  const [files, setFiles] = useState<FileNode[]>([]);
  const [selectedFiles, setSelectedFiles] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);
  const [submitMessage, setSubmitMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Access token state
  const [accessToken, setAccessToken] = useState<string>("");
  const [isUpdatingToken, setIsUpdatingToken] = useState(false);
  const [tokenMessage, setTokenMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);
  const [needsToken, setNeedsToken] = useState(false);

  // Fetch repository details and files on component mount
  useEffect(() => {
    if (repoId) {
      Promise.all([fetchRepository(), fetchFiles()]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [repoId]);

  const fetchRepository = async () => {
    if (!repoId) {
      setError("No repository ID provided.");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const data = await fetchRepositoryDetails(repoId);
      setRepository(data);
    } catch (err) {
      setError("Failed to load repository details. Please try again later.");
      console.error("Error fetching repository:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchFiles = async () => {
    if (!repoId) {
      setFileError("No repository ID provided.");
      return;
    }
    setIsLoading(true);
    setFileError(null);
    setNeedsToken(false);

    try {
      const data = await fetchRepositoryFiles(repoId);
      setFiles(data.files);
    } catch (err) {
      setFileError("Failed to load repository files. Please try again later.");
      console.error("Error fetching files:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const toggleFileSelection = (path: string) => {
    setSelectedFiles((prevSelected) => {
      if (prevSelected.includes(path)) {
        return prevSelected.filter((p) => p !== path);
      } else {
        return [...prevSelected, path];
      }
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    if (!repoId) {
      setSubmitMessage({
        type: "error",
        text: "No repository ID provided.",
      });
      return;
    }

    e.preventDefault();

    if (selectedFiles.length === 0) {
      setSubmitMessage({
        type: "error",
        text: "Please select at least one markdown file to continue",
      });
      return;
    }

    setIsSubmitting(true);
    setSubmitMessage(null);

    try {
      await selectFiles(repoId, selectedFiles);
      setSubmitMessage({
        type: "success",
        text: "Files selected successfully! Redirecting to view markdown content...",
      });
      setTimeout(() => {
        navigate(`/blog?repoId=${repoId}`);
      }, 1500);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      setSubmitMessage({ type: "error", text: message });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle access token update
  const handleTokenSubmit = async (e: React.FormEvent) => {
    if (!repoId) {
      setTokenMessage({
        type: "error",
        text: "No repository ID provided.",
      });
      return;
    }
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
      await updateAccessToken(repoId, accessToken);
      setTokenMessage({
        type: "success",
        text: "Access token updated successfully!",
      });
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

      {error && <ErrorMessage message={error} />}

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

        {fileError && !isLoading && <ErrorMessage message={fileError} />}

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
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="overflow-y-auto max-h-96 border border-gray-200 dark:border-gray-700 rounded p-2">
                <table className="min-w-full">
                  <thead>
                    <tr>
                      <th className="px-4 py-2 text-left text-sm font-medium text-gray-500 dark:text-gray-400">
                        Select
                      </th>
                      <th className="px-4 py-2 text-left text-sm font-medium text-gray-500 dark:text-gray-400">
                        File Path
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {markdownFiles.map((file) => (
                      <tr
                        key={file.path}
                        className="hover:bg-gray-100 dark:hover:bg-gray-700"
                      >
                        <td className="px-4 py-2">
                          <input
                            type="checkbox"
                            checked={selectedFiles.includes(file.path)}
                            onChange={() => toggleFileSelection(file.path)}
                            className="rounded text-blue-500 focus:ring-blue-500"
                          />
                        </td>
                        <td className="px-4 py-2 text-sm">{file.path}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  {selectedFiles.length} file(s) selected
                </span>
                <button
                  type="submit"
                  disabled={isSubmitting || selectedFiles.length === 0}
                  className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
                    isSubmitting || selectedFiles.length === 0
                      ? "opacity-50 cursor-not-allowed"
                      : ""
                  }`}
                >
                  {isSubmitting ? "Processing..." : "Process Selected Files"}
                </button>
              </div>

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
            </form>
          )
        )}
      </div>
    </div>
  );
}

export default RepositoryDetailPage;
