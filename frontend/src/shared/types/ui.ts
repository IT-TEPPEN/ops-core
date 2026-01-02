/**
 * UI type definitions
 */

import { ReactNode } from "react";

/** Common size variants */
export type Size = "sm" | "md" | "lg";

/** Common variant types for components */
export type Variant =
  | "primary"
  | "secondary"
  | "success"
  | "warning"
  | "danger"
  | "info";

/** Status indicator states */
export type StatusType =
  | "active"
  | "inactive"
  | "pending"
  | "error"
  | "success";

/** Base props for form inputs */
export interface BaseInputProps {
  id?: string;
  name: string;
  label?: string;
  error?: string;
  required?: boolean;
  disabled?: boolean;
  className?: string;
}

/** Option type for select inputs */
export interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}
