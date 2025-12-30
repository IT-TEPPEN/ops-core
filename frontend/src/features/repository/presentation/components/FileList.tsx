import { Link } from "react-router-dom";

interface FileNode {
  path: string;
  type: "file" | "dir";
}

interface FileListProps {
  files: FileNode[];
  repositoryId: string;
  isLoading: boolean;
  fileError: string | null;
  needsToken: boolean;
}

export function FileList({
  files,
  repositoryId,
  isLoading,
  fileError,
  needsToken,
}: FileListProps) {
  return (
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

      {isLoading && (
        <p className="text-gray-500">Loading repository files...</p>
      )}

      {!isLoading && !fileError && files.length === 0 ? (
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
                {files.map((file) => (
                  <tr
                    key={file.path}
                    className="hover:bg-gray-100 dark:hover:bg-gray-700"
                  >
                    <td className="px-4 py-2 text-sm">{file.path}</td>
                    <td className="px-4 py-2 text-sm">
                      <Link
                        to={`/repositories/${repositoryId}/files/${encodeURIComponent(
                          file.path
                        )}`}
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
  );
}
