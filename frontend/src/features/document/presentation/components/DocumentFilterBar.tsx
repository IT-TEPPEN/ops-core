interface DocumentFilterBarProps {
  filterType: string;
  searchQuery: string;
  onFilterChange: (filterType: string) => void;
  onSearchChange: (searchQuery: string) => void;
}

export function DocumentFilterBar({
  filterType,
  searchQuery,
  onFilterChange,
  onSearchChange,
}: DocumentFilterBarProps) {
  return (
    <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      <div className="flex flex-col md:flex-row gap-4">
        {/* Search */}
        <div className="flex-1">
          <input
            type="text"
            placeholder="Search by title or tags..."
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        </div>

        {/* Type filter */}
        <div className="md:w-48">
          <select
            value={filterType}
            onChange={(e) => onFilterChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          >
            <option value="all">All Types</option>
            <option value="procedure">Procedures</option>
            <option value="knowledge">Knowledge</option>
          </select>
        </div>
      </div>
    </div>
  );
}
