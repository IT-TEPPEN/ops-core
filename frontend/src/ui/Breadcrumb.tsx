/**
 * Breadcrumb component
 */

import { Link } from "react-router-dom";
import type { BreadcrumbItem } from "../types/ui";

export interface BreadcrumbProps {
  /** Breadcrumb items */
  items: BreadcrumbItem[];
  /** Separator between items */
  separator?: string;
  /** Additional class name */
  className?: string;
}

/**
 * A breadcrumb navigation component
 */
export function Breadcrumb({
  items,
  separator = "/",
  className = "",
}: BreadcrumbProps) {
  return (
    <nav aria-label="Breadcrumb" className={className}>
      <ol className="flex items-center space-x-2 text-sm">
        {items.map((item, index) => {
          const isLast = index === items.length - 1;

          return (
            <li key={index} className="flex items-center">
              {index > 0 && (
                <span
                  className="mx-2 text-gray-400 dark:text-gray-500"
                  aria-hidden="true"
                >
                  {separator}
                </span>
              )}

              {item.href && !isLast ? (
                <Link
                  to={item.href}
                  className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
                >
                  {item.label}
                </Link>
              ) : (
                <span
                  className={
                    isLast
                      ? "text-gray-900 dark:text-gray-100 font-medium"
                      : "text-gray-500 dark:text-gray-400"
                  }
                  aria-current={isLast ? "page" : undefined}
                >
                  {item.label}
                </span>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
