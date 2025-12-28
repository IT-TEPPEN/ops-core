import { Link } from "react-router-dom";

interface RouteErrorProps {
  message: string;
  fallbackPath: string;
  fallbackLabel?: string;
}

export function RouteError({
  message,
  fallbackPath,
  fallbackLabel = "Back",
}: RouteErrorProps) {
  return (
    <div className="space-y-4">
      <div className="p-4 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
        {message}
      </div>
      <Link
        to={fallbackPath}
        className="inline-block px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
      >
        {fallbackLabel}
      </Link>
    </div>
  );
}
