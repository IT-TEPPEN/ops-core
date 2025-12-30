import { useState, ReactNode } from "react";
import { copyToClipboard } from "../utils/copyToClipboard";
import { extractTextContent } from "../utils/extractTextContent";

interface CodeBlockProps {
  children: ReactNode;
  className?: string;
}

export function CodeBlock({ children, className }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    const code = extractTextContent(children);
    await copyToClipboard(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="code-block-wrapper">
      <pre className={className}>{children}</pre>
      <button
        onClick={handleCopy}
        className="code-block-copy-btn"
        aria-label="Copy code to clipboard"
      >
        {copied ? "Copied!" : "Copy"}
      </button>
    </div>
  );
}
