import { Link } from "react-router-dom";
import { useRepositoryManagementAdapter } from "@/features/repository/hooks";
import { useQuery } from "@tanstack/react-query";
import type {
  DocumentKnowledgeMeta,
  DocumentProcedureMeta,
} from "@/features/repository/types";
import { Page } from "@/shared/types/Page";
import { useEffect, useState } from "react";
import { MarkdownProcessor } from "@/features/markdown/processor";
import "@/features/markdown/markdown.css";

function ProcedureMetaComponent(props: { meta: DocumentProcedureMeta }) {
  return (
    <div className="mb-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-2">Procedure Metadata</h2>
      <dl className="space-y-2">
        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">
            Title
          </dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.title}
          </dd>
        </div>

        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">
            Owner
          </dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.owner}
          </dd>
        </div>

        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">Tags</dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.tags.join(", ")}
          </dd>
        </div>
      </dl>
    </div>
  );
}

function KnowledgeMetaComponent(props: { meta: DocumentKnowledgeMeta }) {
  return (
    <div className="mb-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-2">Knowledge Metadata</h2>
      <dl className="space-y-2">
        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">
            Title
          </dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.title}
          </dd>
        </div>

        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">
            Owner
          </dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.owner}
          </dd>
        </div>

        <div>
          <dt className="font-medium text-gray-600 dark:text-gray-400">Tags</dt>
          <dd className="text-gray-800 dark:text-gray-200">
            {props.meta.tags.join(", ")}
          </dd>
        </div>
      </dl>
    </div>
  );
}

export const DocumentPreviewPage: Page<"repoId" | "filePath"> = ({
  path: { repoId, filePath },
}) => {
  const [Component, setComponent] = useState<React.ReactElement | null>(null);
  const adapter = useRepositoryManagementAdapter();
  const query = useQuery({
    queryKey: ["repositories", repoId, "files", filePath],
    queryFn: () => adapter.getFileContent(repoId!, filePath!),
  });

  useEffect(() => {
    if (!query.isLoading && !query.error && query.data) {
      MarkdownProcessor.process(query.data.getContent()).then(
        (file: { result: unknown }) => {
          setComponent(file.result as React.ReactElement);
        }
      );
    }
  }, [query]);

  if (query.isLoading) {
    return (
      <div className="text-center p-8">
        <p className="text-gray-500">Loading content...</p>
      </div>
    );
  }

  if (query.error) {
    return (
      <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
        {(query.error as Error).message}
      </div>
    );
  }

  if (!query.data) {
    return (
      <div className="p-8 bg-white dark:bg-gray-800 rounded-lg shadow-md">
        <p className="text-gray-500">
          No markdown content available for this repository.
        </p>
        <p className="text-gray-500 mt-2">
          Make sure you have selected markdown files from the repository
          management page.
        </p>
      </div>
    );
  }

  const meta = query.data.getMeta();

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Documentation</h1>
        <Link
          to={`/repositories/${repoId}`}
          className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600"
        >
          Back to Repository
        </Link>
      </div>

      <div className="flex flex-col md:flex-row gap-6">
        <div className="w-full md:w-1/4 bg-white dark:bg-gray-800 rounded-lg shadow-md">
          {meta.type === "procedure" && <ProcedureMetaComponent meta={meta} />}
          {meta.type === "knowledge" && <KnowledgeMetaComponent meta={meta} />}
        </div>

        <div className="p-6 md:p-8 bg-white dark:bg-gray-800 rounded-lg shadow-md">
          <article className="markdown-content max-w-none">{Component}</article>
        </div>
      </div>
    </div>
  );
};
