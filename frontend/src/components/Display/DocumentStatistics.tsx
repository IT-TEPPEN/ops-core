import React from "react";
import type { DocumentStatistics as DocumentStatsType } from "../../types/domain";
import { formatDateTime } from "../../utils/date";

export interface DocumentStatisticsProps {
  statistics: DocumentStatsType;
  isLoading?: boolean;
}

/**
 * DocumentStatistics Component
 * Displays statistics for a document
 */
export const DocumentStatistics: React.FC<DocumentStatisticsProps> = ({
  statistics,
  isLoading = false,
}) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const stats = [
    {
      label: "Total Views",
      value: statistics.total_views.toLocaleString(),
      icon: "📊",
    },
    {
      label: "Unique Viewers",
      value: statistics.unique_viewers.toLocaleString(),
      icon: "👥",
    },
    {
      label: "Last Viewed",
      value: formatDateTime(statistics.last_viewed_at),
      icon: "🕒",
    },
    {
      label: "Avg. Duration",
      value: `${statistics.average_view_duration}s`,
      icon: "⏱️",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat) => (
        <div
          key={stat.label}
          className="bg-white rounded-lg shadow p-6 border border-gray-200"
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">{stat.label}</p>
              <p className="mt-2 text-3xl font-semibold text-gray-900">
                {stat.value}
              </p>
            </div>
            <div className="text-4xl">{stat.icon}</div>
          </div>
        </div>
      ))}
    </div>
  );
};
