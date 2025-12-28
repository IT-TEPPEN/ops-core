import { useReducer, useEffect, useCallback } from "react";
import { DocumentListItem } from "@/shared/types/domain";

interface DocumentListState {
  documents: DocumentListItem[];
  isLoading: boolean;
  error: string | null;
  filterType: string;
  searchQuery: string;
}

type DocumentListAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; documents: DocumentListItem[] }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "SET_FILTER_TYPE"; filterType: string }
  | { type: "SET_SEARCH_QUERY"; searchQuery: string };

function documentListReducer(
  state: DocumentListState,
  action: DocumentListAction
): DocumentListState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, error: null };

    case "FETCH_SUCCESS":
      return {
        ...state,
        documents: action.documents,
        isLoading: false,
      };

    case "FETCH_ERROR":
      return {
        ...state,
        isLoading: false,
        error: action.error,
        documents: [],
      };

    case "SET_FILTER_TYPE":
      return { ...state, filterType: action.filterType };

    case "SET_SEARCH_QUERY":
      return { ...state, searchQuery: action.searchQuery };

    default:
      return state;
  }
}

const initialState: DocumentListState = {
  documents: [],
  isLoading: true,
  error: null,
  filterType: "all",
  searchQuery: "",
};

export function useDocumentList() {
  const [state, dispatch] = useReducer(documentListReducer, initialState);

  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  const fetchDocuments = useCallback(async () => {
    dispatch({ type: "FETCH_START" });

    try {
      const response = await fetch(`${apiUrl}/documents`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      dispatch({
        type: "FETCH_SUCCESS",
        documents: data.documents || [],
      });
    } catch (err) {
      dispatch({
        type: "FETCH_ERROR",
        error: "Failed to load documents. Please try again later.",
      });
      console.error("Error fetching documents:", err);
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchDocuments();
  }, [fetchDocuments]);

  const setFilterType = useCallback((filterType: string) => {
    dispatch({ type: "SET_FILTER_TYPE", filterType });
  }, []);

  const setSearchQuery = useCallback((searchQuery: string) => {
    dispatch({ type: "SET_SEARCH_QUERY", searchQuery });
  }, []);

  // Filter documents
  const filteredDocuments = state.documents.filter((doc) => {
    const matchesType =
      state.filterType === "all" || doc.doc_type === state.filterType;
    const matchesSearch =
      state.searchQuery === "" ||
      doc.title.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
      doc.tags.some((tag) =>
        tag.toLowerCase().includes(state.searchQuery.toLowerCase())
      );
    return matchesType && matchesSearch;
  });

  return {
    documents: state.documents,
    filteredDocuments,
    isLoading: state.isLoading,
    error: state.error,
    filterType: state.filterType,
    searchQuery: state.searchQuery,
    setFilterType,
    setSearchQuery,
  };
}
