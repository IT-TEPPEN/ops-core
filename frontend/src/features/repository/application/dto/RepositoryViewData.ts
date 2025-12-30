/**
 * Repository view data for display purposes.
 * Plain object interface (not a class) following ADR 0018/0019.
 */
export interface RepositoryViewData {
  id: string;
  name: string;
  url: string;
  createdAt: Date;
  updatedAt: Date;
}
