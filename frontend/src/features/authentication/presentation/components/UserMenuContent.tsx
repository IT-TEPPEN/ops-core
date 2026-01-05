import { Link } from "react-router-dom";

interface UserMenuContentProps {
  onSignOut: () => void;
  onClose: () => void;
}

export function UserMenuContent({ onSignOut, onClose }: UserMenuContentProps) {
  const handleSignOut = () => {
    onClose();
    onSignOut();
  };

  const handleSettingsClick = () => {
    onClose();
  };

  return (
    <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg border border-gray-200 dark:border-gray-700 z-50">
      <div className="py-1">
        <Link
          to="/settings"
          onClick={handleSettingsClick}
          className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition duration-150 ease-in-out"
        >
          User Settings
        </Link>
        <button
          onClick={handleSignOut}
          className="w-full text-left block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition duration-150 ease-in-out"
        >
          Sign Out
        </button>
      </div>
    </div>
  );
}
