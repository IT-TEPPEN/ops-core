import { useState } from "react";
import { RepositoryList } from "../features/repository/components/RepositoryList";

function RepositoriesPage() {
  const [newRepoUrl, setNewRepoUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitMessage, setSubmitMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  // API base URL - directly use the base URL to avoid recalculation
  const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
  const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitMessage(null);

    try {
      const response = await fetch(`${apiUrl}/repositories`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: newRepoUrl }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to register repository");
      }

      setSubmitMessage({
        type: "success",
        text: "Repository registered successfully!",
      });
      setNewRepoUrl("");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      setSubmitMessage({ type: "error", text: message });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Repository Management</h1>

      {/* Registration Form */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Register New Repository</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="repoUrl" className="block text-sm font-medium mb-1">
              Repository URL
            </label>
            <input
              id="repoUrl"
              type="text"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              placeholder="https://github.com/username/repo.git"
              value={newRepoUrl}
              onChange={(e) => setNewRepoUrl(e.target.value)}
              required
            />
          </div>
          <button
            type="submit"
            disabled={isSubmitting}
            className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
              isSubmitting ? "opacity-50 cursor-not-allowed" : ""
            }`}
          >
            {isSubmitting ? "Registering..." : "Register Repository"}
          </button>
        </form>

        {submitMessage && (
          <div
            className={`mt-4 p-3 rounded ${
              submitMessage.type === "success"
                ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
            }`}
          >
            {submitMessage.text}
          </div>
        )}
      </div>

      <RepositoryList />
    </div>
  );
}

export default RepositoriesPage;
