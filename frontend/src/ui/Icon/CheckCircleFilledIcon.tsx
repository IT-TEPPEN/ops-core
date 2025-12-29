import type { IconProps } from "./types";

/**
 * CheckCircleFilled icon - Filled check circle for selected state
 */
export function CheckCircleFilledIcon({
  className = "",
  size = 24,
}: IconProps) {
  return (
    <svg
      className={className}
      fill="currentColor"
      viewBox="0 0 20 20"
      width={size}
      height={size}
      aria-label="Selected"
    >
      <path
        fillRule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
        clipRule="evenodd"
      />
    </svg>
  );
}
