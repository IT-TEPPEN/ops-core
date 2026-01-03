import { useSearchParams } from "react-router-dom";
// import { useRepositoryQueryService } from "@/features/repository";
import { useQuery } from "@tanstack/react-query";
import type {
  DocumentKnowledgeMeta,
  DocumentProcedureMeta,
  DocumentVariable,
  // DocumentVariable,
} from "@/features/repository/types";
import { Page } from "@/shared/types/Page";
import { useEffect, useState, useCallback } from "react";
import { MarkdownProcessor } from "@/features/markdown/processor";
// import { FileCommitHistory } from "@/features/repository/presentation/components/FileCommitHistory";
import { VariableForm } from "@/features/common/components";
import { substituteVariables } from "@/shared/utils/variableSubstitution";
// import { useDocumentRegistration } from "@/features/document/hooks/useDocumentRegistration";
// import { DocumentRegistrationDialog } from "@/features/document/components/DocumentRegistrationDialog";
// import { useNotifications } from "@/features/notification";
import { useOAuthQueryService } from "@/features/oauth";
import type { VariableDefinition } from "@/shared/types/domain";
import { LoadingSpinner } from "@/ui";

function documentVariableToDefinition(
  docVar: DocumentVariable
): VariableDefinition {
  return {
    name: docVar.name,
    label: docVar.label,
    description: docVar.description ?? undefined,
    type: docVar.type,
    required: docVar.required,
    default_value: docVar.defaultValue ?? undefined,
  };
}

interface ProcedureMetaComponentProps {
  meta: DocumentProcedureMeta;
  variableValues: Record<string, string | number | boolean>;
  onVariableChange: (name: string, value: string | number | boolean) => void;
  onValidate?: () => Promise<boolean>;
}

function ProcedureMetaComponent(props: ProcedureMetaComponentProps) {
  return (
    <div className="mb-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">Procedure Metadata</h2>
      <dl className="space-y-8">
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

      {!!props.meta.variables && props.meta.variables.length > 0 && (
        <div className="mt-8">
          <VariableForm
            variables={props.meta.variables.map(documentVariableToDefinition)}
            values={props.variableValues}
            onChange={props.onVariableChange}
            onValidate={props.onValidate}
          />
        </div>
      )}
    </div>
  );
}

function KnowledgeMetaComponent(props: { meta: DocumentKnowledgeMeta }) {
  return (
    <div className="mb-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-2">Knowledge Metadata</h2>
      <dl className="space-y-8">
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
export const DocumentPreviewPage: Page<
  never,
  "connection_id" | "repository_full_name" | "path"
> = ({ query: { connection_id, repository_full_name, path } }) => {
  const [Component, setComponent] = useState<React.ReactElement | null>(null);
  const [variableValues, setVariableValues] = useState<
    Record<string, string | number | boolean>
  >({});
  const [searchParams] = useSearchParams();
  const commit = searchParams.get("commit") || undefined;
  const [_, setIsDialogOpen] = useState(false);
  const queryService = useOAuthQueryService();
  // const navigate = useNavigate();
  // const { registerDocument, isLoading: isRegistering } =
  //   useDocumentRegistration();
  // const { actions: notificationActions } = useNotifications();

  const query = useQuery({
    queryKey: [
      "connections",
      connection_id,
      "repositories",
      repository_full_name,
      "files",
      path,
      "commits",
      commit,
    ],
    queryFn: () =>
      queryService.getFileContent(
        connection_id!,
        repository_full_name,
        path,
        commit
      ),
  });

  const handleVariableChange = useCallback(
    (name: string, value: string | number | boolean) => {
      setVariableValues((prev) => ({
        ...prev,
        [name]: value,
      }));
    },
    []
  );

  const handleValidate = useCallback(async (): Promise<boolean> => {
    if (!query.data?.meta || query.data.meta.type !== "procedure") return true;

    const variables = query.data.meta.variables;
    if (!variables || variables.length === 0) return true;

    let isValid = true;

    for (const variable of variables) {
      if (variable.required) {
        const value = variableValues[variable.name];
        if (value === undefined || value === null || value === "") {
          isValid = false;
        }
      }
    }

    return isValid;
  }, [query.data, variableValues]);

  // Initialize variable values when document loads
  useEffect(() => {
    if (query.data?.meta.type === "procedure" && query.data.meta.variables) {
      const initialValues: Record<string, string | number | boolean> = {};
      query.data.meta.variables.forEach((v) => {
        initialValues[v.name] = v.defaultValue ?? "";
      });
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setVariableValues(initialValues);
    }
  }, [query.data?.meta]);

  // Process content with variable substitution
  useEffect(() => {
    if (!query.isLoading && !query.error && query.data) {
      const substituted = substituteVariables(
        query.data.content,
        variableValues
      );

      MarkdownProcessor.process(substituted).then(
        (file: { result: unknown }) => {
          setComponent(file.result as React.ReactElement);
        }
      );
    }
  }, [query.data, query.isLoading, query.error, variableValues]);

  // const handleCommitSelect = (commitHash: string) => {
  //   setSearchParams({ commit: commitHash });
  // };

  // const handleRegisterDocument = useCallback(
  //   async (options: {
  //     accessScope: "public" | "private";
  //     isAutoUpdate: boolean;
  //   }) => {
  //     if (!query.data) return;

  //     try {
  //       const result = await registerDocument(
  //         repoId,
  //         filePath!,
  //         options,
  //         commit
  //       );
  //       notificationActions.push({
  //         title: "Success",
  //         message: "Document registered successfully",
  //         type: "success",
  //       });
  //       setIsDialogOpen(false);
  //       navigate(`/documents/${result.id}`);
  //     } catch (error) {
  //       notificationActions.push({
  //         title: "Error",
  //         message:
  //           error instanceof Error
  //             ? error.message
  //             : "Failed to register document",
  //         type: "error",
  //       });
  //     }
  //   },
  //   [
  //     registerDocument,
  //     navigate,
  //     notificationActions,
  //     query.data,
  //     commit,
  //     repoId,
  //     filePath,
  //   ]
  // );

  if (query.isLoading) {
    return <LoadingSpinner message="Loading..." />;
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

  const meta = query.data.meta;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Documentation</h1>
          {query.data?.commitHash && (
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Viewing commit: {query.data.commitHash.substring(0, 7)}
            </p>
          )}
        </div>
        <div className="flex gap-3">
          <button
            onClick={() => setIsDialogOpen(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Register as Document
          </button>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-6">
        <div className="w-full md:w-1/4 space-y-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
            {meta.type === "procedure" && (
              <ProcedureMetaComponent
                meta={meta}
                variableValues={variableValues}
                onVariableChange={handleVariableChange}
                onValidate={handleValidate}
              />
            )}
            {meta.type === "knowledge" && (
              <KnowledgeMetaComponent meta={meta} />
            )}
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
            {/* <FileCommitHistory
              repoId={repoId!}
              filePath={filePath!}
              currentCommit={query.data?.commitHash}
              onCommitSelect={handleCommitSelect}
            /> */}
          </div>
        </div>

        <div className="flex-1 p-6 md:p-8 bg-white dark:bg-gray-800 rounded-lg shadow-md">
          <article className="markdown-content max-w-none">{Component}</article>
        </div>
      </div>

      {/* <DocumentRegistrationDialog
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
        onConfirm={handleRegisterDocument}
        isLoading={isRegistering}
      /> */}
    </div>
  );
};
