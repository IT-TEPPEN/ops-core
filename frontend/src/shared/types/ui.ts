/**
 * UI type definitions
 */

import { ReactNode } from "react";

/** Common size variants */
export type Size = "sm" | "md" | "lg";

/** Common variant types for components */
export type Variant = "primary" | "secondary" | "success" | "warning" | "danger" | "info";

/** Status indicator states */
export type StatusType = "active" | "inactive" | "pending" | "error" | "success";

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

/** Tab item definition */
export interface TabItem {
  id: string;
  label: string;
  content: ReactNode;
  disabled?: boolean;
}

/** Breadcrumb item definition */
export interface BreadcrumbItem {
  label: string;
  href?: string;
}

/** Modal props */
export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  size?: Size;
}

/** Toast notification type */
export interface Toast {
  id: string;
  type: Variant;
  message: string;
  duration?: number;
}

/** Column definition for data tables */
export interface TableColumn<T> {
  key: keyof T | string;
  header: string;
  render?: (item: T) => ReactNode;
  sortable?: boolean;
  width?: string;
}

/** Pagination state */
export interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
}

/** Form field state */
export interface FormFieldState {
  value: unknown;
  touched: boolean;
  error?: string;
}

/** Form state */
export interface FormState<T> {
  values: T;
  errors: Partial<Record<keyof T, string>>;
  touched: Partial<Record<keyof T, boolean>>;
  isSubmitting: boolean;
  isValid: boolean;
}
