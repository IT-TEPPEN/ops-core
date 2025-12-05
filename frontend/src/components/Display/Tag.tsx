/**
 * Tag component
 */

import type { ReactNode } from "react";

export interface TagProps {
  /** Tag content */
  children: ReactNode;
  /** Color (can be any valid CSS color) */
  color?: string;
  /** Remove handler (makes tag removable) */
  onRemove?: () => void;
  /** Click handler */
  onClick?: () => void;
  /** Additional class name */
  className?: string;
}

const defaultColors = [
  { bg: "bg-blue-100 dark:bg-blue-900", text: "text-blue-700 dark:text-blue-200" },
  { bg: "bg-green-100 dark:bg-green-900", text: "text-green-700 dark:text-green-200" },
  { bg: "bg-purple-100 dark:bg-purple-900", text: "text-purple-700 dark:text-purple-200" },
  { bg: "bg-pink-100 dark:bg-pink-900", text: "text-pink-700 dark:text-pink-200" },
  { bg: "bg-orange-100 dark:bg-orange-900", text: "text-orange-700 dark:text-orange-200" },
];

/**
 * A tag component for displaying categories or keywords
 */
export function Tag({
  children,
  color,
  onRemove,
  onClick,
  className = "",
}: TagProps) {
  // Generate a consistent color based on content using a simple hash function
  const getColorIndex = (str: string): number => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32bit integer
    }
    return Math.abs(hash) % defaultColors.length;
  };

  const colorIndex = typeof children === "string" 
    ? getColorIndex(children)
    : 0;
  const defaultColor = defaultColors[colorIndex];

  const hasCustomColor = !!color;
  const customStyle = hasCustomColor ? { backgroundColor: color } : undefined;

  return (
    <span
      onClick={onClick}
      style={customStyle}
      className={`
        inline-flex items-center gap-1
        px-2 py-0.5 rounded text-sm font-medium
        ${!hasCustomColor ? `${defaultColor.bg} ${defaultColor.text}` : "text-white"}
        ${onClick ? "cursor-pointer hover:opacity-80" : ""}
        ${className}
      `}
    >
      {children}
      {onRemove && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onRemove();
          }}
          className="ml-1 hover:opacity-70 focus:outline-none"
          aria-label="Remove tag"
        >
          <svg
            className="w-3 h-3"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clipRule="evenodd"
            />
          </svg>
        </button>
      )}
    </span>
  );
}
