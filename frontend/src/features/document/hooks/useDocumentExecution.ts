import { useState, useEffect } from "react";
import { Document } from "@/shared/types/domain";
import { substituteVariables } from "@/shared/utils/variableSubstitution";

interface UseDocumentExecutionProps {
  docId: string | undefined;
  recordId?: string | undefined;
  initialVariableValues?: Record<string, string | number | boolean>;
}

export function useDocumentExecution(
  propsOrDocId: UseDocumentExecutionProps | string | undefined
) {
  // Support both object props and simple string docId
  const props: UseDocumentExecutionProps =
    typeof propsOrDocId === "string" || propsOrDocId === undefined
      ? { docId: propsOrDocId }
      : propsOrDocId;

  const { docId, recordId, initialVariableValues } = props;
  const [document, setDocument] = useState<Document | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [variableValues, setVariableValues] = useState<
    Record<string, string | number | boolean>
  >({});
  const [processedContent, setProcessedContent] = useState<string>("");
  const [executionTitle, setExecutionTitle] = useState<string>("");

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  // Fetch document on component mount
  useEffect(() => {
    if (!docId) return;

    const fetchDocument = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const response = await fetch(`${apiUrl}/documents/${docId}`);
        if (!response.ok) {
          if (response.status === 404) {
            setError("Document not found");
          } else {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return;
        }
        const data = await response.json();
        setDocument(data);

        // Set default title
        if (!executionTitle) {
          setExecutionTitle(
            `Execution of ${
              data.current_version?.title || "Document"
            } - ${new Date().toLocaleString()}`
          );
        }
      } catch (err) {
        setError("Failed to load document. Please try again later.");
        console.error("Error fetching document:", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDocument();
  }, [docId, apiUrl]);

  // Initialize variable values
  useEffect(() => {
    if (initialVariableValues) {
      setVariableValues(initialVariableValues);
    } else if (document?.current_version?.variables && !recordId) {
      const initialValues: Record<string, string | number | boolean> = {};
      document.current_version.variables.forEach((v) => {
        if (v.default_value !== undefined && v.default_value !== null) {
          initialValues[v.name] = v.default_value;
        } else {
          switch (v.type) {
            case "number":
              initialValues[v.name] = 0;
              break;
            case "boolean":
              initialValues[v.name] = false;
              break;
            default:
              initialValues[v.name] = "";
          }
        }
      });
      setVariableValues(initialValues);
    }
  }, [document, recordId, initialVariableValues]);

  // Process content with variable substitution
  useEffect(() => {
    if (document?.current_version?.content) {
      const substituted = substituteVariables(
        document.current_version.content,
        variableValues
      );
      setProcessedContent(substituted);
    }
  }, [document, variableValues]);

  const handleVariableChange = (
    name: string,
    value: string | number | boolean
  ) => {
    setVariableValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  return {
    document,
    isLoading,
    error,
    variableValues,
    processedContent,
    executionTitle,
    setExecutionTitle,
    handleVariableChange,
  };
}
