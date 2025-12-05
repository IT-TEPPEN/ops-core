/**
 * Pagination component
 */

import type { Size } from "../../types/ui";

export interface PaginationProps {
  /** Current page (1-indexed) */
  currentPage: number;
  /** Total number of pages */
  totalPages: number;
  /** Callback when page changes */
  onPageChange: (page: number) => void;
  /** Maximum number of page buttons to show */
  maxVisible?: number;
  /** Size variant */
  size?: Size;
  /** Show first/last buttons */
  showFirstLast?: boolean;
  /** Additional class name */
  className?: string;
}

const sizeClasses: Record<Size, { button: string; text: string }> = {
  sm: { button: "px-2 py-1", text: "text-sm" },
  md: { button: "px-3 py-1.5", text: "text-base" },
  lg: { button: "px-4 py-2", text: "text-lg" },
};

/**
 * A pagination component for navigating through pages
 */
export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  maxVisible = 7,
  size = "md",
  showFirstLast = true,
  className = "",
}: PaginationProps) {
  const sizes = sizeClasses[size];

  const getPageNumbers = (): (number | "...")[] => {
    if (totalPages <= maxVisible) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }

    const halfVisible = Math.floor(maxVisible / 2);
    const pages: (number | "...")[] = [];

    if (currentPage <= halfVisible + 1) {
      for (let i = 1; i <= maxVisible - 2; i++) {
        pages.push(i);
      }
      pages.push("...");
      pages.push(totalPages);
    } else if (currentPage >= totalPages - halfVisible) {
      pages.push(1);
      pages.push("...");
      for (let i = totalPages - maxVisible + 3; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);
      pages.push("...");
      for (let i = currentPage - 1; i <= currentPage + 1; i++) {
        pages.push(i);
      }
      pages.push("...");
      pages.push(totalPages);
    }

    return pages;
  };

  const pageNumbers = getPageNumbers();

  const buttonBaseClass = `
    ${sizes.button} ${sizes.text}
    rounded border border-gray-300 dark:border-gray-600
    bg-white dark:bg-gray-800
    text-gray-700 dark:text-gray-300
    hover:bg-gray-50 dark:hover:bg-gray-700
    disabled:opacity-50 disabled:cursor-not-allowed
    transition-colors
  `;

  const activeButtonClass = `
    ${sizes.button} ${sizes.text}
    rounded border border-blue-500
    bg-blue-500 text-white
  `;

  return (
    <nav className={`flex items-center gap-1 ${className}`} aria-label="Pagination">
      {/* First page button */}
      {showFirstLast && (
        <button
          onClick={() => onPageChange(1)}
          disabled={currentPage === 1}
          className={buttonBaseClass}
          aria-label="First page"
        >
          «
        </button>
      )}

      {/* Previous page button */}
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className={buttonBaseClass}
        aria-label="Previous page"
      >
        ‹
      </button>

      {/* Page numbers */}
      {pageNumbers.map((page, index) => {
        if (page === "...") {
          return (
            <span
              key={`ellipsis-${index}`}
              className={`${sizes.button} ${sizes.text} text-gray-500`}
            >
              ...
            </span>
          );
        }

        const isActive = page === currentPage;
        return (
          <button
            key={page}
            onClick={() => onPageChange(page)}
            className={isActive ? activeButtonClass : buttonBaseClass}
            aria-current={isActive ? "page" : undefined}
          >
            {page}
          </button>
        );
      })}

      {/* Next page button */}
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className={buttonBaseClass}
        aria-label="Next page"
      >
        ›
      </button>

      {/* Last page button */}
      {showFirstLast && (
        <button
          onClick={() => onPageChange(totalPages)}
          disabled={currentPage === totalPages}
          className={buttonBaseClass}
          aria-label="Last page"
        >
          »
        </button>
      )}
    </nav>
  );
}
