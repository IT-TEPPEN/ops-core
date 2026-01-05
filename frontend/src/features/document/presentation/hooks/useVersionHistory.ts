import { useReducer, useEffect } from "react";
import type { VersionHistoryItem } from "../../application";
import { useDocumentQueryService, useDocumentCommandService } from "../contexts";

type VersionHistoryState = {
  versions: VersionHistoryItem[];
  isLoading: boolean;
  error: string | null;
  selectedVersions: [number | null, number | null];
  actionMessage: { type: "success" | "error"; text: string } | null;
  isProcessing: boolean;
};

type VersionHistoryAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; versions: VersionHistoryItem[] }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "TOGGLE_VERSION"; versionNumber: number }
  | { type: "ACTION_START" }
  | { type: "ACTION_SUCCESS"; message: string }
  | { type: "ACTION_ERROR"; error: string }
  | { type: "CLEAR_MESSAGE" };

const initialState: VersionHistoryState = {
  versions: [],
  isLoading: true,
  error: null,
  selectedVersions: [null, null],
  actionMessage: null,
  isProcessing: false,
};

function versionHistoryReducer(
  state: VersionHistoryState,
  action: VersionHistoryAction
): VersionHistoryState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, error: null };
    case "FETCH_SUCCESS":
      return { ...state, isLoading: false, versions: action.versions };
    case "FETCH_ERROR":
      return { ...state, isLoading: false, error: action.error };
    case "TOGGLE_VERSION": {
      const prev = state.selectedVersions;
      const versionNumber = action.versionNumber;
      let newSelection: [number | null, number | null] = [null, null];

      if (prev[0] === versionNumber) {
        newSelection = [prev[1], null];
      } else if (prev[1] === versionNumber) {
        newSelection = [prev[0], null];
      } else if (prev[0] === null) {
        newSelection = [versionNumber, null];
      } else if (prev[1] === null) {
        newSelection = [prev[0], versionNumber];
      } else {
        newSelection = [prev[1], versionNumber];
      }

      return { ...state, selectedVersions: newSelection };
    }
    case "ACTION_START":
      return { ...state, isProcessing: true, actionMessage: null };
    case "ACTION_SUCCESS":
      return {
        ...state,
        isProcessing: false,
        actionMessage: { type: "success", text: action.message },
      };
    case "ACTION_ERROR":
      return {
        ...state,
        isProcessing: false,
        actionMessage: { type: "error", text: action.error },
      };
    case "CLEAR_MESSAGE":
      return { ...state, actionMessage: null };
    default:
      return state;
  }
}

/**
 * Hook for managing document version history.
 * Refactored to use DocumentQueryService and DocumentCommandService following ADR 0019.
 */
export function useVersionHistory(docId: string | undefined) {
  const [state, dispatch] = useReducer(versionHistoryReducer, initialState);
  const queryService = useDocumentQueryService();
  const commandService = useDocumentCommandService();

  const fetchVersions = async () => {
    if (!docId) return;

    dispatch({ type: "FETCH_START" });

    try {
      const versions = await queryService.getVersionHistory({ docId });
      // Sort versions by version number descending (newest first)
      const sortedVersions = [...versions].sort(
        (a, b) => b.versionNumber - a.versionNumber
      );
      dispatch({ type: "FETCH_SUCCESS", versions: sortedVersions });
    } catch (err) {
      if (err instanceof Error && err.message.includes("404")) {
        dispatch({ type: "FETCH_ERROR", error: "Document not found" });
      } else {
        dispatch({
          type: "FETCH_ERROR",
          error: "Failed to load version history. Please try again later.",
        });
      }
      console.error("Error fetching versions:", err);
    }
  };

  useEffect(() => {
    if (docId) {
      fetchVersions();
    }
  }, [docId]);

  const handleRollback = async (versionNumber: number) => {
    if (
      !window.confirm(
        `Are you sure you want to rollback to version ${versionNumber}?`
      )
    ) {
      return;
    }

    dispatch({ type: "ACTION_START" });

    try {
      // TODO: Implement rollback via CommandService
      // For now, throw an error to indicate it's not implemented
      throw new Error("Rollback functionality not yet implemented in service");
    } catch (err) {
      dispatch({
        type: "ACTION_ERROR",
        error: err instanceof Error ? err.message : "Failed to rollback",
      });
    }
  };

  const handlePublish = async (versionNumber: number) => {
    if (!docId) return;

    dispatch({ type: "ACTION_START" });

    try {
      await commandService.publishVersion({ docId, versionNumber });

      dispatch({
        type: "ACTION_SUCCESS",
        message: `Successfully published version ${versionNumber}`,
      });

      // Refresh versions
      await fetchVersions();
    } catch (err) {
      dispatch({
        type: "ACTION_ERROR",
        error: err instanceof Error ? err.message : "Failed to publish",
      });
    }
  };

  const toggleVersionSelection = (versionNumber: number) => {
    dispatch({ type: "TOGGLE_VERSION", versionNumber });
  };

  return {
    ...state,
    handleRollback,
    handlePublish,
    toggleVersionSelection,
  };
}
