import React from "react";
import type { PopularDocument } from "../../types/domain";
import { formatDateTime } from "../../utils/date";

interface PopularDocumentsListProps {
  items: PopularDocument[];
  isLoading?: boolean;
}

/**
 * PopularDocumentsList Component
 * Displays a ranking of popular documents
 */
export const PopularDocumentsList: React.FC<PopularDocumentsListProps> = ({
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
        No popular documents found.
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <h3 className="text-lg font-semibold text-gray-900">
          Popular Documents
        </h3>
      </div>
      <div className="divide-y divide-gray-200">
        {items.map((item, index) => (
          <div
            key={item.document_id}
            className="px-6 py-4 hover:bg-gray-50 flex items-center justify-between"
          >
            <div className="flex items-center space-x-4">
              <div className="flex-shrink-0">
                <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-blue-100 text-blue-800 font-bold">
                  {index + 1}
                </span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-900">
                  Document {item.document_id.substring(0, 8)}...
                </p>
                <p className="text-xs text-gray-500">
                  Last viewed: {formatDateTime(item.last_viewed_at)}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-6 text-sm">
              <div className="text-right">
                <p className="font-medium text-gray-900">
                  {item.total_views.toLocaleString()}
                </p>
                <p className="text-xs text-gray-500">views</p>
              </div>
              <div className="text-right">
                <p className="font-medium text-gray-900">
                  {item.unique_viewers.toLocaleString()}
                </p>
                <p className="text-xs text-gray-500">viewers</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
