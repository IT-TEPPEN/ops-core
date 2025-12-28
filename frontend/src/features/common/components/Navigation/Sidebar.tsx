/**
 * Sidebar component
 */

import { Link, useLocation } from "react-router-dom";
import type { ReactNode } from "react";

export interface SidebarItem {
  /** Item ID */
  id: string;
  /** Display label */
  label: string;
  /** Navigation href */
  href?: string;
  /** Icon element */
  icon?: ReactNode;
  /** Whether the item is disabled */
  disabled?: boolean;
  /** Click handler (for items without href) */
  onClick?: () => void;
  /** Nested items */
  children?: SidebarItem[];
}

export interface SidebarProps {
  /** Sidebar items */
  items: SidebarItem[];
  /** Header content */
  header?: ReactNode;
  /** Footer content */
  footer?: ReactNode;
  /** Additional class name */
  className?: string;
}

/**
 * A sidebar navigation component
 */
export function Sidebar({
  items,
  header,
  footer,
  className = "",
}: SidebarProps) {
  const location = useLocation();

  const renderItem = (item: SidebarItem, depth = 0) => {
    const isActive = item.href === location.pathname;
    const hasChildren = item.children && item.children.length > 0;

    const itemContent = (
      <>
        {item.icon && (
          <span className="flex-shrink-0 w-5 h-5">{item.icon}</span>
        )}
        <span className="flex-1 truncate">{item.label}</span>
      </>
    );

    const itemClasses = `
      flex items-center gap-3 px-3 py-2 rounded-md
      text-sm font-medium
      transition-colors
      ${depth > 0 ? "ml-6" : ""}
      ${item.disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}
      ${
        isActive
          ? "bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200"
          : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
      }
    `;

    return (
      <li key={item.id}>
        {item.href && !item.disabled ? (
          <Link to={item.href} className={itemClasses}>
            {itemContent}
          </Link>
        ) : (
          <button
            onClick={item.disabled ? undefined : item.onClick}
            className={`w-full text-left ${itemClasses}`}
            disabled={item.disabled}
          >
            {itemContent}
          </button>
        )}
        {hasChildren && (
          <ul className="mt-1 space-y-1">
            {item.children!.map((child) => renderItem(child, depth + 1))}
          </ul>
        )}
      </li>
    );
  };

  return (
    <nav className={`flex flex-col h-full ${className}`}>
      {header && (
        <div className="flex-shrink-0 px-4 py-4 border-b border-gray-200 dark:border-gray-700">
          {header}
        </div>
      )}

      <ul className="flex-1 overflow-y-auto px-2 py-4 space-y-1">
        {items.map((item) => renderItem(item))}
      </ul>

      {footer && (
        <div className="flex-shrink-0 px-4 py-4 border-t border-gray-200 dark:border-gray-700">
          {footer}
        </div>
      )}
    </nav>
  );
}
