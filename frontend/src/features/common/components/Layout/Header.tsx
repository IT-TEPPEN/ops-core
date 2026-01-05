import { Link } from "react-router-dom";
import { UserMenu } from "@/features/authentication/presentation/components";

export function Header() {
  return (
    <nav className="bg-white dark:bg-gray-800 shadow-md h-full">
      <div className="max-w-5xl mx-auto px-4 h-full">
        <div className="flex justify-between items-center h-full">
          <div className="flex-1 flex justify-center">
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
                  to="/documents"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Documents
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
          <div className="flex items-center">
            <UserMenu />
          </div>
        </div>
      </div>
    </nav>
  );
}
