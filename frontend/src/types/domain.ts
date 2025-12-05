/**
 * Domain model type definitions
 */

/** Repository entity */
export interface Repository {
  id: string;
  name: string;
  url: string;
  createdAt: string;
  updatedAt: string;
}

/** File node in repository */
export interface FileNode {
  path: string;
  type: "file" | "dir";
}

/** Document entity */
export interface Document {
  id: string;
  title: string;
  content: string;
  repositoryId: string;
  filePath: string;
  createdAt: string;
  updatedAt: string;
}

/** Document metadata */
export interface DocumentMetadata {
  title?: string;
  author?: string;
  date?: string;
  tags?: string[];
  [key: string]: unknown;
}

/** User entity */
export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
}

/** User roles */
export type UserRole = "admin" | "editor" | "viewer";

/** Execution record */
export interface ExecutionRecord {
  id: string;
  documentId: string;
  executedBy: string;
  executedAt: string;
  status: ExecutionStatus;
  notes?: string;
}

/** Execution status */
export type ExecutionStatus = "pending" | "in_progress" | "completed" | "failed";

/** Evidence entry */
export interface Evidence {
  id: string;
  executionRecordId: string;
  type: EvidenceType;
  content: string;
  createdAt: string;
}

/** Evidence types */
export type EvidenceType = "screenshot" | "log" | "output" | "note";
