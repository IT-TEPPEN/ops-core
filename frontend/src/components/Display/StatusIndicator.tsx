/**
 * Status indicator component
 */

import type { StatusType, Size } from "../../types/ui";

export interface StatusIndicatorProps {
  /** Status type */
  status: StatusType;
  /** Label text */
  label?: string;
  /** Size */
  size?: Size;
  /** Show pulse animation for active status */
  pulse?: boolean;
  /** Additional class name */
  className?: string;
}

const statusColors: Record<StatusType, { dot: string; text: string }> = {
  active: {
    dot: "bg-green-500",
    text: "text-green-700 dark:text-green-300",
  },
  inactive: {
    dot: "bg-gray-400",
    text: "text-gray-600 dark:text-gray-400",
  },
  pending: {
    dot: "bg-yellow-500",
    text: "text-yellow-700 dark:text-yellow-300",
  },
  error: {
    dot: "bg-red-500",
    text: "text-red-700 dark:text-red-300",
  },
  success: {
    dot: "bg-green-500",
    text: "text-green-700 dark:text-green-300",
  },
};

const statusLabels: Record<StatusType, string> = {
  active: "Active",
  inactive: "Inactive",
  pending: "Pending",
  error: "Error",
  success: "Success",
};

const sizeClasses: Record<Size, { dot: string; text: string }> = {
  sm: { dot: "w-1.5 h-1.5", text: "text-xs" },
  md: { dot: "w-2 h-2", text: "text-sm" },
  lg: { dot: "w-2.5 h-2.5", text: "text-base" },
};

/**
 * A status indicator component with dot and optional label
 */
export function StatusIndicator({
  status,
  label,
  size = "md",
  pulse = false,
  className = "",
}: StatusIndicatorProps) {
  const colors = statusColors[status];
  const sizes = sizeClasses[size];
  const displayLabel = label ?? statusLabels[status];

  return (
    <span className={`inline-flex items-center gap-2 ${className}`}>
      <span className="relative flex">
        {pulse && status === "active" && (
          <span
            className={`
              absolute inline-flex h-full w-full
              rounded-full ${colors.dot} opacity-75 animate-ping
            `}
          />
        )}
        <span
          className={`
            relative inline-flex rounded-full
            ${sizes.dot} ${colors.dot}
          `}
        />
      </span>
      {displayLabel && (
        <span className={`font-medium ${sizes.text} ${colors.text}`}>
          {displayLabel}
        </span>
      )}
    </span>
  );
}
