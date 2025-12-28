import { Link } from "react-router-dom";
import { useRepositoryDetail } from "@/features/repository/hooks/useRepositoryDetail";
import { AccessTokenForm } from "@/features/repository/components/AccessTokenForm";
import { FileList } from "@/features/repository/components/FileList";
import { Page } from "@/shared/types/Page";

export const RepositoryDetailPage: Page<"repoId"> = ({ path: { repoId } }) => {
  const {
    repository,
    isLoading,
    error,
    fileError,
    needsToken,
    accessToken,
    isUpdatingToken,
    tokenMessage,
    handleTokenSubmit,
    setAccessToken,
    markdownFiles,
  } = useRepositoryDetail(repoId);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Repository Details</h1>
        <Link
          to="/repositories"
          className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600"
        >
          Back to Repositories
        </Link>
      </div>

      {isLoading && !repository && (
        <p className="text-gray-500">Loading repository information...</p>
      )}

      {error && (
        <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
      )}

      {repository && (
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-2">{repository.name}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            <span className="font-medium">URL:</span> {repository.url}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            <span className="font-medium">Registered on:</span>{" "}
            {new Date(repository.createdAt).toLocaleString()}
          </p>
        </div>
      )}

      {repository && (needsToken || fileError) && (
        <AccessTokenForm
          accessToken={accessToken}
          isUpdatingToken={isUpdatingToken}
          tokenMessage={tokenMessage}
          onTokenChange={setAccessToken}
          onSubmit={handleTokenSubmit}
        />
      )}

      <FileList
        files={markdownFiles}
        repositoryId={repoId}
        isLoading={isLoading}
        fileError={fileError}
        needsToken={needsToken}
      />
    </div>
  );
};
