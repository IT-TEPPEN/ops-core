import { useReducer, useEffect, useCallback } from "react";
import type { DocumentViewData } from "../../application";
import { useDocumentQueryService } from "../contexts";
import { substituteVariables } from "@/shared/utils/variableSubstitution";

interface VariableDefinition {
  name: string;
  default_value?: any;
}

interface DocumentViewState {
  document: DocumentViewData | null;
  isLoading: boolean;
  error: string | null;
  variableValues: Record<string, any>;
  processedContent: string;
}

type DocumentViewAction =
  | { type: "FETCH_START" }
  | {
      type: "FETCH_SUCCESS";
      payload: DocumentViewData;
      variables: VariableDefinition[];
    }
  | { type: "FETCH_ERROR"; payload: string }
  | { type: "SET_VARIABLE"; name: string; value: any }
  | { type: "PROCESS_CONTENT"; content: string };

function documentViewReducer(
  state: DocumentViewState,
  action: DocumentViewAction
): DocumentViewState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, error: null };

    case "FETCH_SUCCESS":
      const initialValues: Record<string, any> = {};
      action.variables?.forEach((v: VariableDefinition) => {
        initialValues[v.name] = v.default_value ?? "";
      });
      return {
        ...state,
        document: action.payload,
        isLoading: false,
        variableValues: initialValues,
      };

    case "FETCH_ERROR":
      return {
        ...state,
        isLoading: false,
        error: action.payload,
        document: null,
      };

    case "SET_VARIABLE":
      return {
        ...state,
        variableValues: {
          ...state.variableValues,
          [action.name]: action.value,
        },
      };

    case "PROCESS_CONTENT":
      return {
        ...state,
        processedContent: action.content,
      };

    default:
      return state;
  }
}

const initialState: DocumentViewState = {
  document: null,
  isLoading: true,
  error: null,
  variableValues: {},
  processedContent: "",
};

/**
 * Hook for viewing document details with variable substitution.
 * Refactored to use DocumentQueryService following ADR 0019.
 */
export function useDocumentView(docId: string | undefined) {
  const [state, dispatch] = useReducer(documentViewReducer, initialState);
  const queryService = useDocumentQueryService();

  // Fetch document
  useEffect(() => {
    if (!docId) return;

    const fetchDocument = async () => {
      dispatch({ type: "FETCH_START" });

      try {
        const document = await queryService.getById({ docId });
        // TODO: Parse variables from document metadata or frontmatter
        const variables: VariableDefinition[] = [];
        dispatch({ type: "FETCH_SUCCESS", payload: document, variables });
      } catch (err) {
        if (err instanceof Error && err.message.includes("404")) {
          dispatch({ type: "FETCH_ERROR", payload: "Document not found" });
        } else {
          dispatch({
            type: "FETCH_ERROR",
            payload: "Failed to load document. Please try again later.",
          });
        }
        console.error("Error fetching document:", err);
      }
    };

    fetchDocument();
  }, [docId, queryService]);

  // Process content when document or variables change
  useEffect(() => {
    if (state.document?.currentVersion?.content) {
      const substituted = substituteVariables(
        state.document.currentVersion.content,
        state.variableValues
      );
      dispatch({ type: "PROCESS_CONTENT", content: substituted });
    }
  }, [state.document, state.variableValues]);

  const handleVariableChange = useCallback((name: string, value: any) => {
    dispatch({ type: "SET_VARIABLE", name, value });
  }, []);

  const handleValidate = useCallback(async (): Promise<boolean> => {
    // TODO: Implement validation via service
    // For now, return true
    return true;
  }, []);

  return {
    document: state.document,
    isLoading: state.isLoading,
    error: state.error,
    variableValues: state.variableValues,
    processedContent: state.processedContent,
    handleVariableChange,
    handleValidate,
  };
}
