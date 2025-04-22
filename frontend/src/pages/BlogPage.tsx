import { useState, useEffect } from "react";
import { useSearchParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";

function BlogPage() {
  const [searchParams] = useSearchParams();
  const repoId = searchParams.get("repoId");

  const [markdownContent, setMarkdownContent] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api";

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
      const response = await fetch(`${apiUrl}/repositories/${repoId}/markdown`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setMarkdownContent(data.content);
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

      {error && (
        <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
      )}

      {!isLoading && !error && markdownContent && (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
          <div className="p-6 md:p-8">
            <article className="prose lg:prose-xl dark:prose-invert max-w-none">
              <ReactMarkdown>{markdownContent}</ReactMarkdown>
            </article>
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
