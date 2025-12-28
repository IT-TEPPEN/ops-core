/**
 * Format utility functions
 */

/**
 * Truncate a string to a maximum length
 */
export const truncate = (str: string, maxLength: number, suffix = "..."): string => {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength - suffix.length) + suffix;
};

/**
 * Capitalize the first letter of a string
 */
export const capitalize = (str: string): string => {
  if (!str) return "";
  return str.charAt(0).toUpperCase() + str.slice(1);
};

/**
 * Convert a string to title case
 */
export const toTitleCase = (str: string): string => {
  return str
    .toLowerCase()
    .split(" ")
    .map((word) => capitalize(word))
    .join(" ");
};

/**
 * Convert camelCase to kebab-case
 */
export const toKebabCase = (str: string): string => {
  return str.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
};

/**
 * Convert kebab-case to camelCase
 */
export const toCamelCase = (str: string): string => {
  return str.replace(/-([a-z])/g, (_, letter) => letter.toUpperCase());
};

/**
 * Format a number with thousands separators
 */
export const formatNumber = (num: number, locale = "ja-JP"): string => {
  return new Intl.NumberFormat(locale).format(num);
};

/**
 * Format bytes to human readable string
 */
export const formatBytes = (bytes: number, decimals = 2): string => {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(decimals)) + " " + sizes[i];
};

/**
 * Format a percentage
 */
export const formatPercentage = (value: number, decimals = 1): string => {
  return `${(value * 100).toFixed(decimals)}%`;
};

/**
 * Format currency
 */
export const formatCurrency = (
  amount: number,
  currency = "JPY",
  locale = "ja-JP"
): string => {
  return new Intl.NumberFormat(locale, {
    style: "currency",
    currency,
  }).format(amount);
};

/**
 * Pluralize a word based on count
 */
export const pluralize = (count: number, singular: string, plural?: string): string => {
  if (count === 1) return singular;
  return plural || `${singular}s`;
};

/**
 * Join strings with proper grammar (e.g., "A, B, and C")
 */
export const joinWithAnd = (items: string[], separator = ", ", finalWord = "and"): string => {
  if (items.length === 0) return "";
  if (items.length === 1) return items[0];
  if (items.length === 2) return items.join(` ${finalWord} `);
  return `${items.slice(0, -1).join(separator)}${separator}${finalWord} ${items[items.length - 1]}`;
};

/**
 * Sanitize a string for use in URLs (slug)
 */
export const toSlug = (str: string): string => {
  return str
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, "")
    .replace(/[\s_-]+/g, "-")
    .replace(/^-+|-+$/g, "");
};

/**
 * Generate initials from a name
 */
export const getInitials = (name: string, maxLength = 2): string => {
  return name
    .split(" ")
    .map((part) => part.charAt(0))
    .join("")
    .toUpperCase()
    .slice(0, maxLength);
};
