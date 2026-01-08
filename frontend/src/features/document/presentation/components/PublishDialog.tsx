import { LoadingSpinner } from "@/ui";
import { useState } from "react";

interface PublishDialogProps {
  filePath: string;
  onClose: () => void;
  onConfirm: (options: { accessScope: "public" | "private" }) => void;
}

const radioOptions = [
  {
    value: "public",
    label: "Public",
    description: "Anyone in the organization can view this document",
  },
  {
    value: "private",
    label: "Private",
    description: "Only specific groups can view this document",
  },
] as const;

/**
 * PublishDialog Component
 *
 * ドキュメント公開確認ダイアログ。
 * アクセススコープ（public/private）を選択して公開を確定する。
 *
 * 責任:
 * - 公開設定の選択（アクセススコープ）
 * - 公開の確認/キャンセル
 */
export function PublishDialog({
  filePath,
  onClose,
  onConfirm,
}: PublishDialogProps) {
  const [accessScope, setAccessScope] = useState<"public" | "private">(
    "public"
  );
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!filePath) return null;

  const handleConfirm = async () => {
    setIsSubmitting(true);
    try {
      await onConfirm({ accessScope });
      onClose();
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
        <h2 className="text-xl font-bold mb-4 text-gray-900 dark:text-gray-100">
          Publish Document
        </h2>

        <div className="space-y-4">
          {/* File Information */}
          <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
            <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
              File
            </div>
            <div className="text-sm text-gray-900 dark:text-gray-100 break-all">
              {filePath}
            </div>
          </div>

          {/* Access Scope Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Access Scope
            </label>
            <div className="space-y-2">
              {radioOptions.map((option) => (
                <label
                  key={option.value}
                  className="flex items-start p-3 border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                >
                  <input
                    type="radio"
                    name="accessScope"
                    value={option.value}
                    checked={accessScope === option.value}
                    onChange={() => setAccessScope(option.value)}
                    disabled={isSubmitting}
                    className="mt-1 mr-3"
                  />
                  <div className="flex-1">
                    <div className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {option.label}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400 mt-0.5">
                      {option.description}
                    </div>
                  </div>
                </label>
              ))}
            </div>
          </div>

          {/* Info Note */}
          <div className="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
            <div className="flex gap-2">
              <svg
                className="w-5 h-5 text-blue-600 dark:text-blue-400 shrink-0"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div className="text-xs text-blue-800 dark:text-blue-200">
                The document will be automatically synced with the Git
                repository. Future updates to the file will be reflected in
                OpsCore.
              </div>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-end gap-3 mt-6">
          <button
            onClick={handleClose}
            disabled={isSubmitting}
            className="w-24 px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleConfirm}
            disabled={isSubmitting}
            className="w-24 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2 transition-colors"
          >
            {isSubmitting ? <LoadingSpinner size="sm" /> : "Publish"}
          </button>
        </div>
      </div>
    </div>
  );
}
