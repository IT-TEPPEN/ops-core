import { useState } from "react";

export interface DocumentRegistrationDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (options: {
    accessScope: "public" | "private";
    isAutoUpdate: boolean;
  }) => void;
  isLoading?: boolean;
}

export function DocumentRegistrationDialog({
  isOpen,
  onClose,
  onConfirm,
  isLoading = false,
}: DocumentRegistrationDialogProps) {
  const [accessScope, setAccessScope] = useState<"public" | "private">(
    "public"
  );
  const [isAutoUpdate, setIsAutoUpdate] = useState(true);

  if (!isOpen) return <></>;

  const handleConfirm = () => {
    onConfirm({ accessScope, isAutoUpdate });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
        <h2 className="text-xl font-bold mb-4 text-gray-900 dark:text-gray-100">
          Register as Document
        </h2>

        <div className="space-y-4">
          {/* Access Scope */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Access Scope
            </label>
            <div className="space-y-2">
              <label className="flex items-center">
                <input
                  type="radio"
                  name="accessScope"
                  value="public"
                  checked={accessScope === "public"}
                  onChange={() => setAccessScope("public")}
                  disabled={isLoading}
                  className="mr-2"
                />
                <span className="text-gray-800 dark:text-gray-200">
                  Public - Anyone can view this document
                </span>
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  name="accessScope"
                  value="private"
                  checked={accessScope === "private"}
                  onChange={() => setAccessScope("private")}
                  disabled={isLoading}
                  className="mr-2"
                />
                <span className="text-gray-800 dark:text-gray-200">
                  Private - Restricted access
                </span>
              </label>
            </div>
          </div>

          {/* Auto Update */}
          <div>
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={isAutoUpdate}
                onChange={(e) => setIsAutoUpdate(e.target.checked)}
                disabled={isLoading}
                className="mr-2"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">
                Enable auto-update (automatically update when repository
                changes)
              </span>
            </label>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-end gap-3 mt-6">
          <button
            onClick={onClose}
            disabled={isLoading}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Cancel
          </button>
          <button
            onClick={handleConfirm}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            {isLoading && (
              <svg
                className="animate-spin h-4 w-4 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            )}
            {isLoading ? "Registering..." : "Register"}
          </button>
        </div>
      </div>
    </div>
  );
}
