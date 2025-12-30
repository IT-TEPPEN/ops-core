import { ReactNode } from "react";

/**
 * React要素からテキストコンテンツを抽出する関数
 */
export function extractTextContent(node: ReactNode): string {
  if (typeof node === "string" || typeof node === "number") {
    return String(node);
  }

  if (Array.isArray(node)) {
    return node.map(extractTextContent).join("");
  }

  if (
    node &&
    typeof node === "object" &&
    "props" in node &&
    node.props &&
    typeof node.props === "object" &&
    "children" in node.props
  ) {
    return extractTextContent(node.props.children as ReactNode);
  }

  return "";
}
