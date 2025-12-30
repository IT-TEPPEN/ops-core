import { Link } from "react-router-dom";
import { ExecutionRecord } from "@/shared/types/domain";

interface ExecutionHeaderProps {
  docId: string;
  executionRecord: ExecutionRecord | null;
  executionTitle: string;
  executionNotes: string;
  isCreating: boolean;
  isSaving: boolean;
  onTitleChange: (title: string) => void;
  onNotesChange: (notes: string) => void;
  onUpdateTitle: () => void;
  onUpdateNotes: () => void;
  onStartExecution: () => void;
  onComplete: () => void;
  onFail: () => void;
}

export function ExecutionHeader({
  docId,
  executionRecord,
  executionTitle,
  executionNotes,
  isCreating,
  isSaving,
  onTitleChange,
  onNotesChange,
  onUpdateTitle,
  onUpdateNotes,
  onStartExecution,
  onComplete,
  onFail,
}: ExecutionHeaderProps) {
  return (
    <div className="mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-4">
          <Link
            to={`/documents/${docId}`}
            className="text-blue-500 hover:text-blue-700"
          >
            ← Back to Document
          </Link>
          <h1 className="text-2xl font-bold">Execute Procedure</h1>
        </div>

        {executionRecord && (
          <div className="flex gap-2">
            {executionRecord.status === "in_progress" && (
              <>
                <button
                  onClick={onComplete}
                  disabled={isSaving}
                  className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:opacity-50"
                >
                  {isSaving ? "Saving..." : "Complete"}
                </button>
                <button
                  onClick={onFail}
                  disabled={isSaving}
                  className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
                >
                  {isSaving ? "Saving..." : "Mark as Failed"}
                </button>
              </>
            )}
            {executionRecord.status === "completed" && (
              <span className="px-4 py-2 bg-green-100 text-green-800 rounded">
                ✓ Completed
              </span>
            )}
            {executionRecord.status === "failed" && (
              <span className="px-4 py-2 bg-red-100 text-red-800 rounded">
                ✗ Failed
              </span>
            )}
          </div>
        )}
      </div>

      {/* Title and Notes */}
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">
            Execution Title
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              value={executionTitle}
              onChange={(e) => onTitleChange(e.target.value)}
              disabled={!!executionRecord}
              className="flex-1 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
            />
            {executionRecord && (
              <button
                onClick={onUpdateTitle}
                disabled={isSaving}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
              >
                Update
              </button>
            )}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Notes</label>
          <div className="flex gap-2">
            <textarea
              value={executionNotes}
              onChange={(e) => onNotesChange(e.target.value)}
              disabled={!executionRecord}
              rows={3}
              className="flex-1 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
              placeholder="Add notes about this execution..."
            />
            {executionRecord && (
              <button
                onClick={onUpdateNotes}
                disabled={isSaving}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
              >
                Update
              </button>
            )}
          </div>
        </div>
      </div>

      {!executionRecord && (
        <div className="mt-4">
          <button
            onClick={onStartExecution}
            disabled={isCreating}
            className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 font-semibold"
          >
            {isCreating ? "Starting..." : "Start Execution"}
          </button>
        </div>
      )}
    </div>
  );
}
