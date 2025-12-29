import { Link } from "react-router-dom";
import { Document } from "@/shared/types/domain";
import { useEffect, useState } from "react";
import { MarkdownProcessor } from "@/features/markdown";

interface DocumentContentPaneProps {
  document: Document;
  processedContent: string;
}

export function DocumentContentPane({
  document,
  processedContent,
}: DocumentContentPaneProps) {
  const [Component, setComponent] = useState<React.ReactElement | null>(null);

  useEffect(() => {
    MarkdownProcessor.process(processedContent).then(
      (file: { result: unknown }) => {
        setComponent(file.result as React.ReactElement);
      }
    );
  }, [processedContent]);

  const currentVersion = document.current_version;

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h1 className="text-3xl font-bold">
            {currentVersion?.title || "Untitled Document"}
          </h1>
          <Link
            to="/documents"
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Back to List
          </Link>
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Version {currentVersion?.version_number || 1} • by {document.owner}
        </p>
      </div>

      {/* Metadata */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-6">
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
        </div>

        {/* Tags */}
        {currentVersion?.tags && currentVersion.tags.length > 0 && (
          <div className="mt-4">
            <span className="font-medium text-sm">Tags:</span>
            <div className="flex flex-wrap gap-1 mt-1">
              {currentVersion.tags.map((tag: string) => (
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

      {/* Document Content */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow prose dark:prose-invert max-w-none">
        {Component}
      </div>
    </div>
  );
}
