/**
 * Page header component
 */

import type { ReactNode } from "react";

export interface PageHeaderProps {
  /** Page title */
  title: string;
  /** Optional subtitle or description */
  subtitle?: string;
  /** Actions to display on the right side */
  actions?: ReactNode;
  /** Content to display below the title */
  children?: ReactNode;
  /** Additional class name */
  className?: string;
}

/**
 * A page header component with title, subtitle, and actions
 */
export function PageHeader({
  title,
  subtitle,
  actions,
  children,
  className = "",
}: PageHeaderProps) {
  return (
    <header className={`mb-6 ${className}`}>
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 truncate">
            {title}
          </h1>
          {subtitle && (
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {subtitle}
            </p>
          )}
        </div>
        {actions && <div className="flex-shrink-0 flex gap-2">{actions}</div>}
      </div>
      {children && <div className="mt-4">{children}</div>}
    </header>
  );
}
