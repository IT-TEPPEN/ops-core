import { Link } from "react-router-dom";
import { ExecutionRecord } from "../../types/domain";

export interface ExecutionRecordListProps {
  records: ExecutionRecord[];
  documentId: string;
}

export function ExecutionRecordList({
  records,
  documentId,
}: ExecutionRecordListProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "completed":
        return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200";
      case "failed":
        return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200";
      case "in_progress":
        return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200";
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200";
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Execution History</h2>
        <Link
          to={`/documents/${documentId}/execute`}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Start New Execution
        </Link>
      </div>

      {records.length === 0 ? (
        <div className="p-8 text-center bg-gray-50 dark:bg-gray-800 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">
            No execution records yet. Start your first execution to track your
            work.
          </p>
          <Link
            to={`/documents/${documentId}/execute`}
            className="inline-block mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Start Execution
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {records.map((record) => (
            <div
              key={record.id}
              className="p-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="font-medium">{record.title}</h3>
                  <div className="mt-2 flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(record.status)}`}
                    >
                      {record.status}
                    </span>
                    <span>Started: {new Date(record.started_at).toLocaleString()}</span>
                    {record.completed_at && (
                      <span>
                        Completed:{" "}
                        {new Date(record.completed_at).toLocaleString()}
                      </span>
                    )}
                  </div>
                  {record.notes && (
                    <p className="mt-2 text-sm text-gray-600 dark:text-gray-300 line-clamp-2">
                      {record.notes}
                    </p>
                  )}
                  <div className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                    {record.steps.length} step{record.steps.length !== 1 ? "s" : ""} recorded
                  </div>
                </div>
                <Link
                  to={`/documents/${documentId}/execute/${record.id}`}
                  className="ml-4 px-4 py-2 text-sm bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 rounded hover:bg-gray-200 dark:hover:bg-gray-600"
                >
                  View Details
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
