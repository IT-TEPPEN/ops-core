/**
 * Tabs component
 */

import { useState, useCallback } from "react";
import type { TabItem } from "../../types/ui";

export interface TabsProps {
  /** Tab items */
  items: TabItem[];
  /** Active tab ID (controlled) */
  activeTab?: string;
  /** Callback when tab changes */
  onTabChange?: (tabId: string) => void;
  /** Additional class name */
  className?: string;
}

/**
 * A tabs component for switching between content sections
 */
export function Tabs({
  items,
  activeTab: controlledActiveTab,
  onTabChange,
  className = "",
}: TabsProps) {
  const [internalActiveTab, setInternalActiveTab] = useState(items[0]?.id ?? "");

  const activeTab = controlledActiveTab ?? internalActiveTab;

  const handleTabClick = useCallback(
    (tabId: string) => {
      if (controlledActiveTab === undefined) {
        setInternalActiveTab(tabId);
      }
      onTabChange?.(tabId);
    },
    [controlledActiveTab, onTabChange]
  );

  const activeItem = items.find((item) => item.id === activeTab);

  return (
    <div className={className}>
      {/* Tab headers */}
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex -mb-px space-x-8" aria-label="Tabs">
          {items.map((item) => {
            const isActive = item.id === activeTab;
            return (
              <button
                key={item.id}
                onClick={() => !item.disabled && handleTabClick(item.id)}
                disabled={item.disabled}
                className={`
                  py-4 px-1 border-b-2 font-medium text-sm
                  transition-colors
                  ${
                    isActive
                      ? "border-blue-500 text-blue-600 dark:text-blue-400"
                      : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600"
                  }
                  ${item.disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}
                `}
                aria-current={isActive ? "page" : undefined}
              >
                {item.label}
              </button>
            );
          })}
        </nav>
      </div>

      {/* Tab content */}
      <div className="pt-4">
        {activeItem?.content}
      </div>
    </div>
  );
}
