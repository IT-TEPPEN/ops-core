/**
 * Error message component
 */

import type { ReactNode } from "react";
import { XCircleIcon, XIcon } from "@/ui";

export interface ErrorMessageProps {
  /** Error title */
  title?: string;
  /** Error message */
  message: string;
  /** Retry callback */
  onRetry?: () => void;
  /** Dismiss callback */
  onDismiss?: () => void;
  /** Additional actions */
  actions?: ReactNode;
  /** Additional class name */
  className?: string;
}

/**
 * An error message component for displaying errors
 */
export function ErrorMessage({
  title = "Error",
  message,
  onRetry,
  onDismiss,
  actions,
  className = "",
}: ErrorMessageProps) {
  return (
    <div
      className={`
        bg-red-50 dark:bg-red-900/20
        border border-red-200 dark:border-red-800
        rounded-lg p-4
        ${className}
      `}
      role="alert"
    >
      <div className="flex">
        <div className="shrink-0">
          <XCircleIcon className="text-red-400" size={20} />
        </div>
        <div className="ml-3 flex-1">
          <h3 className="text-sm font-medium text-red-800 dark:text-red-200">
            {title}
          </h3>
          <p className="mt-1 text-sm text-red-700 dark:text-red-300">
            {message}
          </p>
          {(onRetry || onDismiss || actions) && (
            <div className="mt-4 flex gap-3">
              {onRetry && (
                <button
                  onClick={onRetry}
                  className="text-sm font-medium text-red-800 dark:text-red-200 hover:text-red-600 dark:hover:text-red-100 underline"
                >
                  Retry
                </button>
              )}
              {onDismiss && (
                <button
                  onClick={onDismiss}
                  className="text-sm font-medium text-red-800 dark:text-red-200 hover:text-red-600 dark:hover:text-red-100 underline"
                >
                  Dismiss
                </button>
              )}
              {actions}
            </div>
          )}
        </div>
        {onDismiss && (
          <div className="ml-auto pl-3">
            <button
              onClick={onDismiss}
              className="inline-flex text-red-400 hover:text-red-500 focus:outline-none"
              aria-label="Dismiss"
            >
              <XIcon size={20} />
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
