/**
 * Card component
 */

import type { ReactNode } from "react";

export interface CardProps {
  /** Card title */
  title?: string;
  /** Card subtitle */
  subtitle?: string;
  /** Card content */
  children: ReactNode;
  /** Footer content */
  footer?: ReactNode;
  /** Whether the card is clickable */
  onClick?: () => void;
  /** Additional class name */
  className?: string;
  /** Padding size */
  padding?: "none" | "sm" | "md" | "lg";
}

const paddingClasses = {
  none: "",
  sm: "p-3",
  md: "p-4",
  lg: "p-6",
};

/**
 * A card component for displaying content in a contained box
 */
export function Card({
  title,
  subtitle,
  children,
  footer,
  onClick,
  className = "",
  padding = "md",
}: CardProps) {
  const Component = onClick ? "button" : "div";

  return (
    <Component
      onClick={onClick}
      className={`
        bg-white dark:bg-gray-800
        rounded-lg shadow-md
        overflow-hidden
        ${onClick ? "cursor-pointer hover:shadow-lg transition-shadow text-left w-full" : ""}
        ${className}
      `}
    >
      {(title || subtitle) && (
        <div className={`${paddingClasses[padding]} border-b border-gray-200 dark:border-gray-700`}>
          {title && (
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {title}
            </h3>
          )}
          {subtitle && (
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {subtitle}
            </p>
          )}
        </div>
      )}
      <div className={paddingClasses[padding]}>{children}</div>
      {footer && (
        <div
          className={`${paddingClasses[padding]} bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700`}
        >
          {footer}
        </div>
      )}
    </Component>
  );
}
