import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import rehypeHighlight from "rehype-highlight";
import rehypeReact from "rehype-react";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeAutolinkHeadings from "rehype-autolink-headings";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import remarkToc from "remark-toc";
import { createElement } from "react";
import * as prod from "react/jsx-runtime";
import { CodeBlock } from "./components/CodeBlock";
import { InlineCode } from "./components/InlineCode";
// import { visit } from "unist-util-visit";

// function DebugPlugin() {
//   return transformer;
//   function transformer(tree: any, _file: any) {
//     visit(tree, (node, _index, _parent) => {
//       // コードブロック関連のノードだけをログ出力
//       if (
//         node.type === "element" &&
//         (node.tagName === "pre" || node.tagName === "code")
//       ) {
//         console.log(
//           "Tag:",
//           node.tagName,
//           "Properties:",
//           JSON.stringify(node.properties, null, 2),
//           "Children:",
//           node.children?.length
//         );
//         if (node.tagName === "code" && node.children) {
//           console.log("  First child:", node.children[0]);
//         }
//       }
//     });
//   }
// }

export const MarkdownProcessor = unified()
  .use(remarkParse) // Markdownをパース
  .use(remarkFrontmatter) // Frontmatterをサポート
  .use(remarkGfm) // GitHub Flavored Markdownをサポート
  .use(remarkToc) // 目次生成をサポート
  .use(remarkRehype, { allowDangerousHtml: false }) // MarkdownからHTMLへ変換
  .use(rehypeSlug) // ID付与
  .use(rehypeAutolinkHeadings) // リンク追加
  .use(rehypeHighlight) // シンタックスハイライトをサポート
  // .use(DebugPlugin) // ASTをデバッグ表示（highlightの後）
  .use(rehypeSanitize, {
    ...defaultSchema,
    tagNames: [...(defaultSchema.tagNames || []), "span"],
    attributes: {
      ...defaultSchema.attributes,
      // シンタックスハイライト用のクラスを許可
      code: [["className", "hljs", /^language-/]],
      span: [["className", /^hljs-/, "icon", "icon-link"]],
      pre: [["className", "hljs"]],
      // rehypeAutolinkHeadings用のa属性を許可
      a: [
        ...(defaultSchema.attributes?.a || []),
        ["ariaHidden", "true"],
        ["tabIndex", "-1"],
      ],
    },
  }) // HTMLをサニタイズ
  .use(rehypeReact, {
    ...prod,
    createElement,
    components: {
      pre: CodeBlock, // コードブロック（```で囲まれた部分）にコピー機能を追加
      code: InlineCode, // インラインコード（`で囲まれた部分）にコピー機能を追加
    },
  });
