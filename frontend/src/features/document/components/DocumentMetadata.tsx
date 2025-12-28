import { Document } from "@/shared/types/domain";

interface DocumentMetadataProps {
  document: Document;
}

export function DocumentMetadata({ document }: DocumentMetadataProps) {
  const currentVersion = document.current_version;

  return (
    <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      <div className="flex flex-wrap gap-4 text-sm">
        <div>
          <span className="font-medium">Type:</span>{" "}
          <span
            className={`px-2 py-0.5 rounded-full ${
              currentVersion?.doc_type === "procedure"
                ? "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200"
                : "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
            }`}
          >
            {currentVersion?.doc_type || "unknown"}
          </span>
        </div>
        <div>
          <span className="font-medium">Status:</span>{" "}
          <span
            className={`px-2 py-0.5 rounded-full ${
              document.is_published
                ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
            }`}
          >
            {document.is_published ? "Published" : "Draft"}
          </span>
        </div>
        <div>
          <span className="font-medium">Access:</span>{" "}
          <span className="capitalize">{document.access_scope}</span>
        </div>
        <div>
          <span className="font-medium">Auto Update:</span>{" "}
          {document.is_auto_update ? "Enabled" : "Disabled"}
        </div>
      </div>

      {/* Tags */}
      {currentVersion?.tags && currentVersion.tags.length > 0 && (
        <div className="mt-4">
          <span className="font-medium text-sm">Tags:</span>
          <div className="flex flex-wrap gap-1 mt-1">
            {currentVersion.tags.map((tag) => (
              <span
                key={tag}
                className="px-2 py-0.5 text-xs bg-gray-100 dark:bg-gray-700 rounded"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
