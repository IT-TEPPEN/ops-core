import { Link } from "react-router-dom";
import { useRepositoryList } from "../hooks";

export function RepositoryList() {
  const query = useRepositoryList();

  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">Registered Repositories</h2>

      {query.isLoading && (
        <p className="text-gray-500">Loading repositories...</p>
      )}

      {query.error && (
        <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {query.error.message}
        </div>
      )}

      {!query.isLoading &&
        !query.error &&
        query.data &&
        query.data.pagination.totalItems === 0 && (
          <p className="text-gray-500">No repositories registered yet.</p>
        )}

      {query.data && query.data.pagination.totalItems > 0 && (
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
              {query.data.data.map((repo) => (
                <tr key={repo.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {repo.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {repo.url}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {repo.createdAt.toLocaleString()}
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
  );
}
