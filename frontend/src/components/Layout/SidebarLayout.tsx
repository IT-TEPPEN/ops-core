/**
 * Sidebar layout component
 */

import type { ReactNode } from "react";

export interface SidebarLayoutProps {
  /** Sidebar content */
  sidebar: ReactNode;
  /** Main content */
  children: ReactNode;
  /** Sidebar position */
  sidebarPosition?: "left" | "right";
  /** Sidebar width */
  sidebarWidth?: string;
  /** Whether the sidebar is collapsible */
  collapsible?: boolean;
  /** Whether the sidebar is collapsed (controlled) */
  isCollapsed?: boolean;
  /** Callback when collapse state changes */
  onCollapseChange?: (collapsed: boolean) => void;
  /** Additional class name */
  className?: string;
}

/**
 * A flexible sidebar layout component
 */
export function SidebarLayout({
  sidebar,
  children,
  sidebarPosition = "left",
  sidebarWidth = "w-64",
  collapsible = false,
  isCollapsed = false,
  onCollapseChange,
  className = "",
}: SidebarLayoutProps) {
  const sidebarClasses = `
    ${isCollapsed ? "w-0 overflow-hidden" : sidebarWidth}
    flex-shrink-0
    bg-white dark:bg-gray-800
    border-gray-200 dark:border-gray-700
    transition-all duration-300 ease-in-out
    ${sidebarPosition === "left" ? "border-r" : "border-l"}
  `;

  const sidebarElement = (
    <aside className={sidebarClasses}>
      <div className="h-full overflow-y-auto">
        {sidebar}
      </div>
    </aside>
  );

  return (
    <div className={`flex h-full min-h-0 ${className}`}>
      {sidebarPosition === "left" && sidebarElement}

      <div className="flex-1 min-w-0 flex flex-col relative">
        {collapsible && (
          <button
            onClick={() => onCollapseChange?.(!isCollapsed)}
            className={`
              absolute top-2 z-10
              ${sidebarPosition === "left" ? "left-2" : "right-2"}
              p-1 rounded
              bg-gray-200 dark:bg-gray-700
              hover:bg-gray-300 dark:hover:bg-gray-600
              transition-colors
            `}
            aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
          >
            <svg
              className={`w-5 h-5 transition-transform ${
                isCollapsed
                  ? sidebarPosition === "left"
                    ? ""
                    : "rotate-180"
                  : sidebarPosition === "left"
                    ? "rotate-180"
                    : ""
              }`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>
        )}
        <main className="flex-1 overflow-y-auto bg-gray-50 dark:bg-gray-900">
          {children}
        </main>
      </div>

      {sidebarPosition === "right" && sidebarElement}
    </div>
  );
}
