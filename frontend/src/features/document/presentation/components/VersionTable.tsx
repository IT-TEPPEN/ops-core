import { Link } from "react-router-dom";
import { DocumentVersion } from "@/shared/types/domain";

interface VersionTableProps {
  versions: DocumentVersion[];
  selectedVersions: [number | null, number | null];
  isProcessing: boolean;
  onToggleVersion: (versionNumber: number) => void;
  onRollback: (versionNumber: number) => void;
  onPublish: (versionNumber: number) => void;
  docId: string;
}

export function VersionTable({
  versions,
  selectedVersions,
  isProcessing,
  onToggleVersion,
  onRollback,
  onPublish,
  docId,
}: VersionTableProps) {
  if (versions.length === 0) {
    return (
      <div className="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow">
        <p className="text-gray-500">No versions found</p>
      </div>
    );
  }

  return (
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
                  onChange={() => onToggleVersion(version.version_number)}
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
                      onClick={() => onRollback(version.version_number)}
                      disabled={isProcessing}
                      className="text-orange-600 hover:underline text-sm disabled:opacity-50"
                    >
                      Rollback
                    </button>
                  )}
                  {!version.is_current && version.unpublished_at && (
                    <button
                      onClick={() => onPublish(version.version_number)}
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
  );
}
