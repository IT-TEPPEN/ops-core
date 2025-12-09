import React from "react";
import type { RecentDocument } from "../../types/domain";
import { formatDateTime } from "../../utils/date";

interface RecentViewsListProps {
  items: RecentDocument[];
  isLoading?: boolean;
}

/**
 * RecentViewsList Component
 * Displays a list of recently viewed documents
 */
export const RecentViewsList: React.FC<RecentViewsListProps> = ({
  items,
  isLoading = false,
}) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="text-center p-8 text-gray-500">
        No recently viewed documents found.
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <h3 className="text-lg font-semibold text-gray-900">
          Recently Viewed
        </h3>
      </div>
      <div className="divide-y divide-gray-200">
        {items.map((item) => (
          <div
            key={item.document_id}
            className="px-6 py-4 hover:bg-gray-50 flex items-center justify-between"
          >
            <div>
              <p className="text-sm font-medium text-gray-900">
                Document {item.document_id.substring(0, 8)}...
              </p>
              <p className="text-xs text-gray-500 mt-1">
                {formatDateTime(item.last_viewed_at)}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm font-medium text-gray-900">
                {item.total_views.toLocaleString()}
              </p>
              <p className="text-xs text-gray-500">views</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
