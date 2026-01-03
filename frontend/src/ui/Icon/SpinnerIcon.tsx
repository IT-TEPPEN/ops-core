import type { IconProps } from "./types";

/**
 * Spinner icon - Loading indicator with animation
 */
export function SpinnerIcon({ className = "", size = 24 }: IconProps) {
  return (
    <div
      className={`animate-spin rounded-full border-b-2 ${className}`}
      style={{ height: `${size}px`, width: `${size}px` }}
    ></div>
  );
}
