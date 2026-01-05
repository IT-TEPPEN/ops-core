import { useReducer, useEffect, useCallback } from "react";
import type { DocumentListItem } from "../../application";
import { useDocumentQueryService } from "../contexts";

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

/**
 * Hook for managing document list with filtering and search.
 * Refactored to use DocumentQueryService following ADR 0019.
 */
export function useDocumentList() {
  const [state, dispatch] = useReducer(documentListReducer, initialState);
  const queryService = useDocumentQueryService();

  const fetchDocuments = useCallback(async () => {
    dispatch({ type: "FETCH_START" });

    try {
      const documents = await queryService.list();
      dispatch({
        type: "FETCH_SUCCESS",
        documents,
      });
    } catch (err) {
      dispatch({
        type: "FETCH_ERROR",
        error: "Failed to load documents. Please try again later.",
      });
      console.error("Error fetching documents:", err);
    }
  }, [queryService]);

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
      state.filterType === "all" || doc.docType === state.filterType;
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
