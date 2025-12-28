import { useReducer, useEffect, useCallback } from "react";
import { Document, VariableDefinition } from "@/shared/types/domain";
import { substituteVariables } from "@/shared/utils/variableSubstitution";

interface DocumentViewState {
  document: Document | null;
  isLoading: boolean;
  error: string | null;
  variableValues: Record<string, any>;
  processedContent: string;
}

type DocumentViewAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; payload: Document }
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
      action.payload.current_version?.variables?.forEach(
        (v: VariableDefinition) => {
          initialValues[v.name] = v.default_value ?? "";
        }
      );
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

export function useDocumentView(docId: string | undefined) {
  const [state, dispatch] = useReducer(documentViewReducer, initialState);

  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  // Fetch document
  useEffect(() => {
    if (!docId) return;

    const fetchDocument = async () => {
      dispatch({ type: "FETCH_START" });

      try {
        const response = await fetch(`${apiUrl}/documents/${docId}`);
        if (!response.ok) {
          if (response.status === 404) {
            dispatch({ type: "FETCH_ERROR", payload: "Document not found" });
          } else {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return;
        }
        const data = await response.json();
        dispatch({ type: "FETCH_SUCCESS", payload: data });
      } catch (err) {
        dispatch({
          type: "FETCH_ERROR",
          payload: "Failed to load document. Please try again later.",
        });
        console.error("Error fetching document:", err);
      }
    };

    fetchDocument();
  }, [docId, apiUrl]);

  // Process content when document or variables change
  useEffect(() => {
    if (state.document?.current_version?.content) {
      const substituted = substituteVariables(
        state.document.current_version.content,
        state.variableValues
      );
      dispatch({ type: "PROCESS_CONTENT", content: substituted });
    }
  }, [state.document, state.variableValues]);

  const handleVariableChange = useCallback((name: string, value: any) => {
    dispatch({ type: "SET_VARIABLE", name, value });
  }, []);

  const handleValidate = useCallback(async (): Promise<boolean> => {
    if (!docId) return false;

    try {
      const values = Object.entries(state.variableValues).map(
        ([name, value]) => ({
          name,
          value,
        })
      );

      const response = await fetch(
        `${apiUrl}/documents/${docId}/validate-variables`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ values }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      return result.valid;
    } catch (err) {
      console.error("Error validating variables:", err);
      return false;
    }
  }, [docId, state.variableValues, apiUrl]);

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
