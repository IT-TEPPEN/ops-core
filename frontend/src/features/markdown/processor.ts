import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import rehypeHighlight from "rehype-highlight";
import rehypeReact from "rehype-react";
import remarkFrontmatter from "remark-frontmatter";
// import rehypeStringify from "rehype-stringify";
import remarkGfm from "remark-gfm";
import { createElement } from "react";
import * as prod from "react/jsx-runtime";
import { visit } from "unist-util-visit";

function DebugPlugin() {
  return transformer;
  function transformer(tree: any, _file: any) {
    visit(tree, (node, _index, _parent) => {
      console.log(node);
    });
  }
}

export const MarkdownProcessor = unified()
  .use(remarkParse) // Markdownをパース
  .use(remarkFrontmatter) // Frontmatterをサポート
  .use(DebugPlugin) // ASTをデバッグ表示
  .use(remarkGfm) // GitHub Flavored Markdownをサポート
  .use(remarkRehype) // MarkdownからHTMLへ変換
  .use(rehypeHighlight) // シンタックスハイライトをサポート
  .use(rehypeReact, { ...prod, createElement }); // HTMLをReactコンポーネントに変換
// .use(rehypeReact, {
//   ...prod,
//   // components: {
//   //   // 必要に応じてカスタムコンポーネントを指定可能
//   //   h1: ({ children }: { children: React.ReactNode }) => (
//   //     <h1 className="text-3xl font-bold my-6 border-b border-gray-700 pb-2">
//   //       {children}
//   //     </h1>
//   //   ),
//   //   h2: ({ children }: { children: React.ReactNode }) => (
//   //     <h2 className="text-2xl font-semibold my-5">{children}</h2>
//   //   ),
//   //   h3: ({ children }: { children: React.ReactNode }) => (
//   //     <h3 className="text-xl font-semibold my-4">{children}</h3>
//   //   ),
//   //   strong: ({ children }: { children: React.ReactNode }) => (
//   //     <strong className="font-black text-gray-800">{children}</strong>
//   //   ),
//   //   p: ({ children }: { children: React.ReactNode }) => (
//   //     <p className="my-2 leading-relaxed">{children}</p>
//   //   ),
//   //   a: ({ href, children }: { href?: string; children: React.ReactNode }) => (
//   //     <a href={href} className="text-blue-600 hover:underline">
//   //       {children}
//   //     </a>
//   //   ),
//   // },
// });
