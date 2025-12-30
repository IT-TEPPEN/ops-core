/**
 * Loading component
 */

import { Size } from "@/shared/types/ui";
import { SpinnerIcon } from "@/ui";

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

const sizeMap: Record<Size, number> = {
  sm: 16,
  md: 32,
  lg: 48,
};

const textSizeClasses: Record<Size, string> = {
  sm: "text-sm",
  md: "text-base",
  lg: "text-lg",
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
  const content = (
    <div
      className={`flex flex-col items-center justify-center gap-3 ${className}`}
    >
      <SpinnerIcon className="text-blue-500" size={sizeMap[size]} />
      {text && (
        <p
          className={`text-gray-600 dark:text-gray-400 ${textSizeClasses[size]}`}
        >
          {text}
        </p>
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
