import React from "react";
import { Link } from "react-router-dom";

interface Repository {
  id: string;
  name: string;
  url: string;
  createdAt: string;
}

interface RepositoryTableProps {
  repositories: Repository[];
}

const RepositoryTable: React.FC<RepositoryTableProps> = ({ repositories }) => {
  return (
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
  );
};

export default RepositoryTable;
