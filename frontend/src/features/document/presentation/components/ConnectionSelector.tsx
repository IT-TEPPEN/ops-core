import { UI_Icon_PlusIcon } from "@/ui";
import { Link } from "react-router-dom";

export function ConnectionSelector() {
  return (
    <div className="flex items-center justify-between">
      <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
        Connection
      </h3>
      <Link
        to="/oauth/authorize"
        className="flex items-center gap-1 text-xs px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors cursor-pointer"
      >
        <UI_Icon_PlusIcon size={12} />
        Add
      </Link>
    </div>
  );
}
