import { useQuery } from "@tanstack/react-query";
// import axios from "axios";
import { Link } from "react-router-dom";
import { useRepositoryManagementAdapter } from "../contexts";

// const listRepositories = async (): Promise<Repository[]> => {
//   const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
//   const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

//   return axios
//     .get<{ repositories: Repository[] }>(`${apiUrl}/repositories`)
//     .then((response) => response.data.repositories);
// };

export function RepositoryList() {
  const adapter = useRepositoryManagementAdapter();
  const query = useQuery({
    queryKey: ["repositories"],
    queryFn: () => adapter.listRepositories(),
  });

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
        query.data.pagenation.getTotalItems() === 0 && (
          <p className="text-gray-500">No repositories registered yet.</p>
        )}

      {query.data && query.data.pagenation.getTotalItems() > 0 && (
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
                <tr key={repo.getId()}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {repo.getName()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {repo.getUrl()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {repo.getCreatedAt().toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <Link
                      to={`/repositories/${repo.getId()}`}
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
