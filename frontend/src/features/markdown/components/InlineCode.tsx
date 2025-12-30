import { useState, ReactNode } from "react";
import { copyToClipboard } from "../utils/copyToClipboard";
import { extractTextContent } from "../utils/extractTextContent";

interface InlineCodeProps {
  children: ReactNode;
  className?: string;
}

export function InlineCode({ children, className }: InlineCodeProps) {
  const [copied, setCopied] = useState(false);

  // コードブロック内のcodeタグの場合（language-xxxクラスを持つ）は、
  // 通常のcodeタグとしてレンダリング（コピー機能なし）
  if (className?.includes("language-") || className?.includes("hljs")) {
    return <code className={className}>{children}</code>;
  }

  const handleCopy = async () => {
    const code = extractTextContent(children);
    await copyToClipboard(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <span className="inline-code-wrapper">
      <code className={className}>{children}</code>
      <button
        onClick={handleCopy}
        className="inline-code-copy-btn"
        aria-label="Copy code to clipboard"
      >
        {copied ? "✓" : "📋"}
      </button>
    </span>
  );
}
