import React, { useState, useEffect } from "react";
import { PopularDocumentsList } from "../components/Display/PopularDocumentsList";
import { RecentViewsList } from "../components/Display/RecentViewsList";
import { StatisticsChart } from "../components/Display/StatisticsChart";
import type { PopularDocument, RecentDocument } from "../types/domain";

/**
 * StatisticsPage
 * Displays statistics dashboard with popular and recent documents
 */
const StatisticsPage: React.FC = () => {
  const [popularDocuments, setPopularDocuments] = useState<PopularDocument[]>(
    []
  );
  const [recentDocuments, setRecentDocuments] = useState<RecentDocument[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchStatistics = async () => {
      try {
        setIsLoading(true);
        // TODO: Replace with actual API calls
        // const popularResponse = await fetch('/api/statistics/popular-documents?limit=10');
        // const recentResponse = await fetch('/api/statistics/recent-documents?limit=10');
        // const popularData = await popularResponse.json();
        // const recentData = await recentResponse.json();
        // setPopularDocuments(popularData.items);
        // setRecentDocuments(recentData.items);

        // Mock data for now
        setPopularDocuments([
          {
            document_id: "doc-1",
            total_views: 1234,
            unique_viewers: 567,
            last_viewed_at: new Date().toISOString(),
          },
          {
            document_id: "doc-2",
            total_views: 987,
            unique_viewers: 432,
            last_viewed_at: new Date().toISOString(),
          },
        ]);

        setRecentDocuments([
          {
            document_id: "doc-3",
            last_viewed_at: new Date().toISOString(),
            total_views: 50,
          },
          {
            document_id: "doc-4",
            last_viewed_at: new Date(
              Date.now() - 3600000
            ).toISOString(),
            total_views: 30,
          },
        ]);
      } catch (error) {
        console.error("Failed to fetch statistics:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStatistics();
  }, []);

  const chartData = popularDocuments.map((doc) => ({
    label: `Doc ${doc.document_id.substring(0, 8)}`,
    value: doc.total_views,
  }));

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Statistics</h1>
        <p className="mt-2 text-gray-600">
          View analytics and insights about document usage
        </p>
      </div>

      <div className="space-y-8">
        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <StatisticsChart
            title="Document Views"
            data={chartData}
            isLoading={isLoading}
          />
          <StatisticsChart
            title="View Trends"
            data={[
              { label: "This Week", value: 250 },
              { label: "Last Week", value: 200 },
              { label: "2 Weeks Ago", value: 180 },
              { label: "3 Weeks Ago", value: 150 },
            ]}
            isLoading={isLoading}
          />
        </div>

        {/* Lists Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <PopularDocumentsList items={popularDocuments} isLoading={isLoading} />
          <RecentViewsList items={recentDocuments} isLoading={isLoading} />
        </div>
      </div>
    </div>
  );
};

export default StatisticsPage;
