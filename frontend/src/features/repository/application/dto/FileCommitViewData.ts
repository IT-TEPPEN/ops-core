/**
 * File commit information view data.
 * Represents a single commit in a file's history.
 */
export interface FileCommitViewData {
  commitHash: string;
  message: string;
  author: string;
  authorEmail: string;
  date: Date;
}
