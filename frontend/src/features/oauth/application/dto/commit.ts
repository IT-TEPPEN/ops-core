/**
 * Represents a commit in a file's history.
 * Used for displaying commit history of files in repositories.
 */
export interface Commit {
  /** Commit hash (SHA) */
  hash: string;

  /** Commit message */
  message: string;

  /** Author name */
  author: string;

  /** Author email */
  authorEmail: string;

  /** Commit date */
  date: Date;
}
