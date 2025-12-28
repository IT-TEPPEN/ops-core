/**
 * Date processing utility functions
 */

/**
 * Format a date to ISO string (YYYY-MM-DD)
 */
export const formatDateISO = (date: Date | string): string => {
  const d = typeof date === "string" ? new Date(date) : date;
  return d.toISOString().split("T")[0];
};

/**
 * Format a date to locale string
 */
export const formatDateLocale = (
  date: Date | string,
  locale = "ja-JP",
  options?: Intl.DateTimeFormatOptions
): string => {
  const d = typeof date === "string" ? new Date(date) : date;
  const defaultOptions: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: "long",
    day: "numeric",
    ...options,
  };
  return d.toLocaleDateString(locale, defaultOptions);
};

/**
 * Format a date with time
 */
export const formatDateTime = (
  date: Date | string,
  locale = "ja-JP"
): string => {
  const d = typeof date === "string" ? new Date(date) : date;
  return d.toLocaleString(locale);
};

/**
 * Format a date to relative time (e.g., "3 days ago")
 */
export const formatRelativeTime = (date: Date | string): string => {
  const d = typeof date === "string" ? new Date(date) : date;
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffSecs < 60) {
    return "just now";
  } else if (diffMins < 60) {
    return `${diffMins} minute${diffMins > 1 ? "s" : ""} ago`;
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours > 1 ? "s" : ""} ago`;
  } else if (diffDays < 30) {
    return `${diffDays} day${diffDays > 1 ? "s" : ""} ago`;
  } else {
    return formatDateLocale(d);
  }
};

/**
 * Check if a date is today
 */
export const isToday = (date: Date | string): boolean => {
  const d = typeof date === "string" ? new Date(date) : date;
  const today = new Date();
  return (
    d.getDate() === today.getDate() &&
    d.getMonth() === today.getMonth() &&
    d.getFullYear() === today.getFullYear()
  );
};

/**
 * Check if a date is in the past
 */
export const isPast = (date: Date | string): boolean => {
  const d = typeof date === "string" ? new Date(date) : date;
  return d.getTime() < Date.now();
};

/**
 * Check if a date is in the future
 */
export const isFuture = (date: Date | string): boolean => {
  const d = typeof date === "string" ? new Date(date) : date;
  return d.getTime() > Date.now();
};

/**
 * Add days to a date
 */
export const addDays = (date: Date | string, days: number): Date => {
  const d = typeof date === "string" ? new Date(date) : new Date(date);
  d.setDate(d.getDate() + days);
  return d;
};

/**
 * Get the start of day
 */
export const startOfDay = (date: Date | string): Date => {
  const d = typeof date === "string" ? new Date(date) : new Date(date);
  d.setHours(0, 0, 0, 0);
  return d;
};

/**
 * Get the end of day
 */
export const endOfDay = (date: Date | string): Date => {
  const d = typeof date === "string" ? new Date(date) : new Date(date);
  d.setHours(23, 59, 59, 999);
  return d;
};

/**
 * Parse a date string safely
 */
export const parseDate = (dateString: string): Date | null => {
  const d = new Date(dateString);
  return isNaN(d.getTime()) ? null : d;
};
