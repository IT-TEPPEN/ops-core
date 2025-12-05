/**
 * Loading component
 */

import type { Size } from "../../types/ui";

export interface LoadingProps {
  /** Size variant */
  size?: Size;
  /** Loading text */
  text?: string;
  /** Whether to show full-screen overlay */
  fullScreen?: boolean;
  /** Additional class name */
  className?: string;
}

const sizeClasses: Record<Size, { spinner: string; text: string }> = {
  sm: { spinner: "w-4 h-4", text: "text-sm" },
  md: { spinner: "w-8 h-8", text: "text-base" },
  lg: { spinner: "w-12 h-12", text: "text-lg" },
};

/**
 * A loading spinner component
 */
export function Loading({
  size = "md",
  text,
  fullScreen = false,
  className = "",
}: LoadingProps) {
  const sizes = sizeClasses[size];

  const content = (
    <div className={`flex flex-col items-center justify-center gap-3 ${className}`}>
      <svg
        className={`animate-spin ${sizes.spinner} text-blue-500`}
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        role="status"
        aria-label="Loading"
      >
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
        />
        <path
          className="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        />
      </svg>
      {text && (
        <span className={`${sizes.text} text-gray-500 dark:text-gray-400`}>
          {text}
        </span>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm">
        {content}
      </div>
    );
  }

  return content;
}
