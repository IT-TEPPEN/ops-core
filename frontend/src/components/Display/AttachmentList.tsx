/**
 * AttachmentList component for displaying a list of attachments
 */

import type { AttachmentResponse } from "./AttachmentUploader";

export interface AttachmentListProps {
  /** List of attachments */
  attachments: AttachmentResponse[];
  /** Callback when an attachment is clicked */
  onAttachmentClick?: (attachment: AttachmentResponse) => void;
  /** Callback when delete is clicked */
  onDelete?: (attachmentId: string) => void;
  /** Whether to show delete button */
  showDelete?: boolean;
  /** Additional class name */
  className?: string;
}

/**
 * Formats file size to human-readable string
 */
function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i];
}

/**
 * Formats date to readable string
 */
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

/**
 * Gets icon based on MIME type
 */
function getFileIcon(mimeType: string): string {
  if (mimeType.startsWith("image/")) return "🖼️";
  if (mimeType.startsWith("video/")) return "🎥";
  if (mimeType.startsWith("audio/")) return "🎵";
  if (mimeType.includes("pdf")) return "📄";
  if (mimeType.includes("word") || mimeType.includes("document")) return "📝";
  if (mimeType.includes("excel") || mimeType.includes("spreadsheet")) return "📊";
  if (mimeType.includes("text")) return "📃";
  return "📎";
}

/**
 * Component for displaying a list of attachments
 */
export function AttachmentList({
  attachments,
  onAttachmentClick,
  onDelete,
  showDelete = false,
  className = "",
}: AttachmentListProps) {
  if (attachments.length === 0) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <p className="text-gray-500 dark:text-gray-400">
          No attachments
        </p>
      </div>
    );
  }

  return (
    <div className={`space-y-2 ${className}`}>
      {attachments.map((attachment) => (
        <div
          key={attachment.id}
          className="flex items-center justify-between p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 transition-colors"
        >
          <div
            className="flex items-center flex-1 min-w-0 cursor-pointer"
            onClick={() => onAttachmentClick?.(attachment)}
          >
            <span className="text-2xl mr-3 flex-shrink-0">
              {getFileIcon(attachment.mime_type)}
            </span>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                {attachment.file_name}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {formatFileSize(attachment.file_size)} • {formatDate(attachment.uploaded_at)}
              </p>
            </div>
          </div>
          
          <div className="flex items-center space-x-2 ml-4">
            <a
              href={`/api/v1/attachments/${attachment.id}/download`}
              download
              className="p-2 text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
              title="Download"
              onClick={(e) => e.stopPropagation()}
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                />
              </svg>
            </a>
            
            {showDelete && onDelete && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  if (window.confirm("Are you sure you want to delete this attachment?")) {
                    onDelete(attachment.id);
                  }
                }}
                className="p-2 text-gray-600 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors"
                title="Delete"
              >
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
