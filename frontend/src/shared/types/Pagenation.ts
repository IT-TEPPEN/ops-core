export interface Pagenation {
  getCurrentPage(): number;
  getPreviousPage(): number | null;
  getNextPage(): number | null;
  getTotalPages(): number;
  getParPage(): number;
  getCurrentItems(): number;
  getTotalItems(): number;
}
