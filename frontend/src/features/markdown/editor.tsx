import { useEffect, useState } from "react";
import { MarkdownProcessor } from "./processor";

export const MarkdownTestEditor = () => {
  const [testText, setTestText] = useState("");
  const [resultText, setResultText] = useState<React.ReactElement | null>(null);

  useEffect(() => {
    MarkdownProcessor.process(testText).then((file: { result: unknown }) => {
      setResultText(file.result as React.ReactElement);
    });
  }, [testText]);
  return (
    <div className="grid grid-cols-2 gap-4 w-full h-full">
      <div>
        <textarea
          className="bg-white rounded-2xl w-full h-full p-4 border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-200 focus:border-transparent"
          onChange={(event) => setTestText(event.target.value)}
        />
      </div>
      <div className="w-full p-4 markdown-content">{resultText ?? ""}</div>
    </div>
  );
};
