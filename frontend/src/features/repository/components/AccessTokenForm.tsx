interface AccessTokenFormProps {
  accessToken: string;
  isUpdatingToken: boolean;
  tokenMessage: { type: "success" | "error"; text: string } | null;
  onTokenChange: (token: string) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export function AccessTokenForm({
  accessToken,
  isUpdatingToken,
  tokenMessage,
  onTokenChange,
  onSubmit,
}: AccessTokenFormProps) {
  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">Repository Access Token</h2>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        This repository requires an access token to view files. Please enter a
        valid access token below.
      </p>

      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="accessToken"
            className="block text-sm font-medium mb-1"
          >
            Access Token
          </label>
          <input
            id="accessToken"
            type="password"
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
            placeholder="Enter GitHub personal access token"
            value={accessToken}
            onChange={(e) => onTokenChange(e.target.value)}
            required
          />
          <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
            For GitHub repositories, create a personal access token with 'repo'
            scope.
          </p>
        </div>

        <div>
          <button
            type="submit"
            disabled={isUpdatingToken}
            className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
              isUpdatingToken ? "opacity-50 cursor-not-allowed" : ""
            }`}
          >
            {isUpdatingToken ? "Updating..." : "Update Access Token"}
          </button>
        </div>

        {tokenMessage && (
          <div
            className={`mt-4 p-3 rounded ${
              tokenMessage.type === "success"
                ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
            }`}
          >
            {tokenMessage.text}
          </div>
        )}
      </form>
    </div>
  );
}
