/**
 * AttachmentViewer component for viewing attachment previews
 */

import type { AttachmentResponse } from "./AttachmentUploader";

export interface AttachmentViewerProps {
  /** Attachment to view */
  attachment: AttachmentResponse;
  /** Callback when close is clicked */
  onClose?: () => void;
  /** Additional class name */
  className?: string;
}

/**
 * Component for viewing attachment previews
 */
export function AttachmentViewer({
  attachment,
  onClose,
  className = "",
}: AttachmentViewerProps) {
  const isImage = attachment.mime_type.startsWith("image/");
  const downloadUrl = `/api/v1/attachments/${attachment.id}/download`;

  return (
    <div className={`fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 ${className}`}>
      <div className="relative w-full max-w-4xl max-h-screen p-4">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl overflow-hidden">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
            <div className="flex-1 min-w-0">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 truncate">
                {attachment.file_name}
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {attachment.mime_type}
              </p>
            </div>
            <div className="flex items-center space-x-2 ml-4">
              <a
                href={downloadUrl}
                download
                className="p-2 text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                title="Download"
              >
                <svg
                  className="w-6 h-6"
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
              <button
                onClick={onClose}
                className="p-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
                title="Close"
              >
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="p-4 max-h-[calc(100vh-200px)] overflow-auto">
            {isImage ? (
              <div className="flex justify-center">
                <img
                  src={downloadUrl}
                  alt={attachment.file_name}
                  className="max-w-full h-auto rounded-lg"
                />
              </div>
            ) : (
              <div className="text-center py-12">
                <svg
                  className="w-24 h-24 mx-auto text-gray-400 dark:text-gray-500 mb-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1}
                    d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                  />
                </svg>
                <p className="text-gray-600 dark:text-gray-400 mb-4">
                  Preview not available for this file type
                </p>
                <a
                  href={downloadUrl}
                  download
                  className="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  <svg
                    className="w-5 h-5 mr-2"
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
                  Download File
                </a>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
