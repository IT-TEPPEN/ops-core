import { useState, useEffect } from "react";
import { useSearchParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { fetchRepositoryFiles } from "../api/repositories";
import ErrorMessage from "../ui/ErrorMessage";

interface Metadata {
  [key: string]: any;
}

/**
 * Custom function to parse frontmatter from markdown content
 * This avoids using gray-matter which depends on Node.js Buffer
 */
function parseFrontmatter(markdown: string): {
  content: string;
  data: Metadata;
} {
  const result: { content: string; data: Metadata } = {
    content: markdown,
    data: {},
  };

  // Check if the content starts with frontmatter delimiters (---)
  if (!markdown.startsWith("---")) {
    return result;
  }

  // Find the end of the frontmatter block
  const endOfFrontmatter = markdown.indexOf("---", 3);
  if (endOfFrontmatter === -1) {
    return result;
  }

  // Extract frontmatter text
  const frontmatterText = markdown.substring(3, endOfFrontmatter).trim();

  // Parse frontmatter lines into key-value pairs
  const frontmatterLines = frontmatterText.split("\n");
  frontmatterLines.forEach((line) => {
    const colonIndex = line.indexOf(":");
    if (colonIndex !== -1) {
      const key = line.substring(0, colonIndex).trim();
      const value = line.substring(colonIndex + 1).trim();

      // Try to parse as JSON if possible
      try {
        result.data[key] = JSON.parse(value);
      } catch (e) {
        // If not valid JSON, use the string value
        result.data[key] = value;
      }
    }
  });

  // Return content without frontmatter
  result.content = markdown.substring(endOfFrontmatter + 3).trim();
  return result;
}

function BlogPage() {
  const [searchParams] = useSearchParams();
  const repoId = searchParams.get("repoId");

  const [markdownContent, setMarkdownContent] = useState<string>("");
  const [metadata, setMetadata] = useState<Metadata>({});
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (repoId) {
      fetchMarkdownContent();
    } else {
      setIsLoading(false);
      setError("No repository ID provided. Please select a repository first.");
    }
  }, [repoId]);

  const fetchMarkdownContent = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await fetchRepositoryFiles(repoId);

      // Use custom function to parse the markdown and extract frontmatter
      const { content, data: frontmatter } = parseFrontmatter(data.content);

      // Set the main content without frontmatter
      setMarkdownContent(content);

      // Set the extracted metadata
      setMetadata(frontmatter);
    } catch (err) {
      setError("Failed to load markdown content. Please try again later.");
      console.error("Error fetching markdown content:", err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Documentation</h1>
        <Link
          to="/repositories"
          className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600"
        >
          Back to Repositories
        </Link>
      </div>

      {isLoading && (
        <div className="text-center p-8">
          <p className="text-gray-500">Loading content...</p>
        </div>
      )}

      {error && <ErrorMessage message={error} />}

      {!isLoading && !error && markdownContent && (
        <div className="flex flex-col md:flex-row gap-6">
          {/* Metadata sidebar */}
          {Object.keys(metadata).length > 0 && (
            <div className="w-full md:w-1/4 bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
              <div className="p-4">
                <h2 className="text-lg font-semibold mb-3 pb-2 border-b border-gray-200 dark:border-gray-700">
                  Document Metadata
                </h2>
                <dl className="space-y-2">
                  {Object.entries(metadata).map(([key, value]) => (
                    <div key={key} className="pb-1">
                      <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                        {key}
                      </dt>
                      <dd className="text-sm text-gray-700 dark:text-gray-300 break-words">
                        {typeof value === "string"
                          ? value
                          : JSON.stringify(value)}
                      </dd>
                    </div>
                  ))}
                </dl>
              </div>
            </div>
          )}

          {/* Main content */}
          <div
            className={`w-full ${
              Object.keys(metadata).length > 0 ? "md:w-3/4" : ""
            } bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden`}
          >
            <div className="p-6 md:p-8">
              <article className="prose lg:prose-xl dark:prose-invert max-w-none">
                <ReactMarkdown>{markdownContent}</ReactMarkdown>
              </article>
            </div>
          </div>
        </div>
      )}

      {!isLoading && !error && !markdownContent && (
        <div className="text-center p-8 bg-white dark:bg-gray-800 rounded-lg shadow-md">
          <p className="text-gray-500">
            No markdown content available for this repository.
          </p>
          <p className="text-gray-500 mt-2">
            Make sure you have selected markdown files from the repository
            management page.
          </p>
        </div>
      )}
    </div>
  );
}

export default BlogPage;
