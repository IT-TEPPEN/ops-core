import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { DocumentVersion } from "../types/domain";

interface VersionHistoryResponse {
  document_id: string;
  versions: DocumentVersion[];
}

function DocumentVersionHistoryPage() {
  const { docId } = useParams<{ docId: string }>();

  const [versions, setVersions] = useState<DocumentVersion[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedVersions, setSelectedVersions] = useState<[number | null, number | null]>([null, null]);
  const [actionMessage, setActionMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  // Fetch versions on component mount
  useEffect(() => {
    if (docId) {
      fetchVersions();
    }
  }, [docId]);

  const fetchVersions = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${apiUrl}/documents/${docId}/versions`);
      if (!response.ok) {
        if (response.status === 404) {
          setError("Document not found");
        } else {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return;
      }
      const data: VersionHistoryResponse = await response.json();
      // Sort versions by version number descending (newest first)
      const sortedVersions = [...(data.versions || [])].sort(
        (a, b) => b.version_number - a.version_number
      );
      setVersions(sortedVersions);
    } catch (err) {
      setError("Failed to load version history. Please try again later.");
      console.error("Error fetching versions:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRollback = async (versionNumber: number) => {
    if (!window.confirm(`Are you sure you want to rollback to version ${versionNumber}?`)) {
      return;
    }

    setIsProcessing(true);
    setActionMessage(null);

    try {
      const response = await fetch(
        `${apiUrl}/documents/${docId}/versions/${versionNumber}/rollback`,
        { method: "POST" }
      );

      if (!response.ok) {
        let message = "Failed to rollback";
        try {
          const data = await response.json();
          if (data && typeof data.message === "string") {
            message = data.message;
          }
        } catch (parseError) {
          // Use default message if JSON parsing fails
          console.error('Failed to parse error response:', parseError);
        }
        throw new Error(message);
      }

      setActionMessage({
        type: "success",
        text: `Successfully rolled back to version ${versionNumber}`,
      });

      // Refresh versions
      await fetchVersions();
    } catch (err) {
      setActionMessage({
        type: "error",
        text: err instanceof Error ? err.message : "Failed to rollback",
      });
    } finally {
      setIsProcessing(false);
    }
  };

  const handlePublish = async (versionNumber: number) => {
    setIsProcessing(true);
    setActionMessage(null);

    try {
      const response = await fetch(
        `${apiUrl}/documents/${docId}/versions/${versionNumber}/publish`,
        { method: "POST" }
      );

      if (!response.ok) {
        let message = "Failed to publish";
        try {
          const data = await response.json();
          if (data && typeof data.message === "string") {
            message = data.message;
          }
        } catch {
          // Use default message if JSON parsing fails
        }
        throw new Error(message);
      }

      setActionMessage({
        type: "success",
        text: `Successfully published version ${versionNumber}`,
      });

      // Refresh versions
      await fetchVersions();
    } catch (err) {
      setActionMessage({
        type: "error",
        text: err instanceof Error ? err.message : "Failed to publish",
      });
    } finally {
      setIsProcessing(false);
    }
  };

  const toggleVersionSelection = (versionNumber: number) => {
    setSelectedVersions((prev) => {
      if (prev[0] === versionNumber) {
        return [prev[1], null];
      }
      if (prev[1] === versionNumber) {
        return [prev[0], null];
      }
      if (prev[0] === null) {
        return [versionNumber, null];
      }
      if (prev[1] === null) {
        return [prev[0], versionNumber];
      }
      return [prev[1], versionNumber];
    });
  };

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Loading version history...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        <div className="p-4 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
        <Link
          to="/documents"
          className="inline-block px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
        >
          Back to Documents
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Version History</h1>
        <div className="flex gap-2">
          <Link
            to={`/documents/${docId}`}
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Back to Document
          </Link>
        </div>
      </div>

      {/* Action message */}
      {actionMessage && (
        <div
          className={`p-3 rounded ${
            actionMessage.type === "success"
              ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
              : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
          }`}
        >
          {actionMessage.text}
        </div>
      )}

      {/* Version comparison hint */}
      {selectedVersions[0] !== null && selectedVersions[1] === null && (
        <div className="p-3 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded">
          Select another version to compare
        </div>
      )}

      {/* Version list */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-900">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Compare
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Version
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Title
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Commit
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Published
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {versions.map((version) => (
              <tr
                key={version.id}
                className={`${
                  selectedVersions.includes(version.version_number)
                    ? "bg-blue-50 dark:bg-blue-900/20"
                    : "hover:bg-gray-100 dark:hover:bg-gray-700"
                }`}
              >
                <td className="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    checked={selectedVersions.includes(version.version_number)}
                    onChange={() => toggleVersionSelection(version.version_number)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="font-medium">v{version.version_number}</span>
                </td>
                <td className="px-6 py-4">
                  <div className="text-sm font-medium">{version.title}</div>
                  <div className="text-xs text-gray-500">{version.file_path}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <code className="text-xs bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">
                    {version.commit_hash?.slice(0, 7)}
                  </code>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(version.published_at).toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {version.is_current ? (
                    <span className="px-2 py-1 text-xs bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 rounded-full">
                      Current
                    </span>
                  ) : version.unpublished_at ? (
                    <span className="px-2 py-1 text-xs bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 rounded-full">
                      Unpublished
                    </span>
                  ) : (
                    <span className="px-2 py-1 text-xs bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded-full">
                      Published
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex gap-2">
                    <Link
                      to={`/documents/${docId}/versions/${version.version_number}`}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      View
                    </Link>
                    {!version.is_current && !version.unpublished_at && (
                      <button
                        onClick={() => handleRollback(version.version_number)}
                        disabled={isProcessing}
                        className="text-orange-600 hover:underline text-sm disabled:opacity-50"
                      >
                        Rollback
                      </button>
                    )}
                    {!version.is_current && version.unpublished_at && (
                      <button
                        onClick={() => handlePublish(version.version_number)}
                        disabled={isProcessing}
                        className="text-green-600 hover:underline text-sm disabled:opacity-50"
                      >
                        Publish
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Empty state */}
      {versions.length === 0 && (
        <div className="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow">
          <p className="text-gray-500">No versions found</p>
        </div>
      )}
    </div>
  );
}

export default DocumentVersionHistoryPage;
