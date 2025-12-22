import { Link } from "react-router-dom";

export function Header() {
  return (
    <nav className="bg-white dark:bg-gray-800 shadow-md sticky top-0 z-10">
      <div className="max-w-5xl mx-auto px-4">
        <div className="flex justify-center items-center h-16">
          <ul className="flex space-x-6">
            <li>
              <Link
                to="/"
                className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
              >
                Home
              </Link>
            </li>
            <li>
              <Link
                to="/repositories"
                className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
              >
                Repositories
              </Link>
            </li>
            <li>
              <Link
                to="/documents"
                className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
              >
                Documents
              </Link>
            </li>
            <li>
              <Link
                to="/blog"
                className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
              >
                Documentation
              </Link>
            </li>
            <li>
              <Link
                to="/groups"
                className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
              >
                Groups
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
