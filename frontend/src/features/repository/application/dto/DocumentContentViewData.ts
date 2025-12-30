import { DocumentMeta } from "../../types/repository";

/**
 * Document content view data with parsed frontmatter.
 * Plain object interface replacing DocumentImpl class.
 */
export interface DocumentContentViewData {
  repoId: string;
  filePath: string;
  content: string;
  meta: DocumentMeta;
}
