/**
 * Data table component
 */

import type { ReactNode } from "react";
import type { TableColumn } from "../../types/ui";

export interface DataTableProps<T> {
  /** Data to display */
  data: T[];
  /** Column definitions */
  columns: TableColumn<T>[];
  /** Key field for row identification */
  keyField: keyof T;
  /** Loading state */
  isLoading?: boolean;
  /** Empty state message */
  emptyMessage?: string;
  /** Click handler for rows */
  onRowClick?: (item: T) => void;
  /** Additional class name */
  className?: string;
}

/**
 * A data table component for displaying tabular data
 */
export function DataTable<T>({
  data,
  columns,
  keyField,
  isLoading = false,
  emptyMessage = "No data available",
  onRowClick,
  className = "",
}: DataTableProps<T>) {
  if (isLoading) {
    return (
      <div className="animate-pulse">
        <div className="h-10 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="h-12 bg-gray-100 dark:bg-gray-800 rounded mb-2"
          ></div>
        ))}
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500 dark:text-gray-400">
        {emptyMessage}
      </div>
    );
  }

  const getCellValue = (item: T, column: TableColumn<T>): ReactNode => {
    if (column.render) {
      return column.render(item);
    }
    const value = item[column.key as keyof T];
    return String(value ?? "");
  };

  return (
    <div className={`overflow-x-auto ${className}`}>
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-900">
          <tr>
            {columns.map((column) => (
              <th
                key={String(column.key)}
                scope="col"
                className={`
                  px-6 py-3 text-left text-xs font-medium
                  text-gray-500 dark:text-gray-400 uppercase tracking-wider
                  ${column.width ? column.width : ""}
                `}
              >
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
          {data.map((item) => (
            <tr
              key={String(item[keyField])}
              onClick={() => onRowClick?.(item)}
              className={`
                ${onRowClick ? "cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700" : ""}
                transition-colors
              `}
            >
              {columns.map((column) => (
                <td
                  key={String(column.key)}
                  className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100"
                >
                  {getCellValue(item, column)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
