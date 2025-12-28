/**
 * Three-pane layout component for work evidence display
 */

import type { ReactNode } from "react";

export interface ThreePaneLayoutProps {
  /** Left pane content (navigation/tree) */
  leftPane: ReactNode;
  /** Center pane content (main content) */
  centerPane: ReactNode;
  /** Right pane content (details/metadata) */
  rightPane: ReactNode;
  /** Width of the left pane */
  leftWidth?: string;
  /** Width of the right pane */
  rightWidth?: string;
  /** Additional class name */
  className?: string;
}

/**
 * A three-pane layout component for displaying work evidence
 * with navigation on the left, content in the center, and details on the right
 */
export function ThreePaneLayout({
  leftPane,
  centerPane,
  rightPane,
  leftWidth = "w-64",
  rightWidth = "w-80",
  className = "",
}: ThreePaneLayoutProps) {
  return (
    <div className={`flex h-full min-h-0 ${className}`}>
      {/* Left Pane */}
      <aside
        className={`${leftWidth} flex-shrink-0 border-r border-gray-200 dark:border-gray-700 overflow-y-auto bg-white dark:bg-gray-800`}
      >
        {leftPane}
      </aside>

      {/* Center Pane */}
      <main className="flex-1 min-w-0 overflow-y-auto bg-gray-50 dark:bg-gray-900">
        {centerPane}
      </main>

      {/* Right Pane */}
      <aside
        className={`${rightWidth} flex-shrink-0 border-l border-gray-200 dark:border-gray-700 overflow-y-auto bg-white dark:bg-gray-800`}
      >
        {rightPane}
      </aside>
    </div>
  );
}
