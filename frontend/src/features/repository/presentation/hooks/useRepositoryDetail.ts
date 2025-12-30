import { useReducer, useEffect, useCallback } from "react";
import { useRepositoryQueryService } from "./useRepositoryQueryService";
import { useRepositoryCommandService } from "./useRepositoryCommandService";
import type { FileNodeViewData, RepositoryViewData } from "../../application";
import { AxiosError } from "axios";

type RepositoryDetailState = {
  repository: RepositoryViewData | null;
  files: FileNodeViewData[];
  isLoading: boolean;
  error: string | null;
  fileError: string | null;
  accessToken: string;
  isUpdatingToken: boolean;
  tokenMessage: { type: "success" | "error"; text: string } | null;
  needsToken: boolean;
};

type RepositoryDetailAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_REPOSITORY_SUCCESS"; repository: RepositoryViewData }
  | { type: "FETCH_FILES_SUCCESS"; files: FileNodeViewData[] }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "FETCH_FILE_ERROR"; error: string; needsToken?: boolean }
  | { type: "SET_ACCESS_TOKEN"; token: string }
  | { type: "UPDATE_TOKEN_START" }
  | { type: "UPDATE_TOKEN_SUCCESS"; message: string }
  | { type: "UPDATE_TOKEN_ERROR"; error: string }
  | { type: "CLEAR_TOKEN_MESSAGE" };

const initialState: RepositoryDetailState = {
  repository: null,
  files: [],
  isLoading: true,
  error: null,
  fileError: null,
  accessToken: "",
  isUpdatingToken: false,
  tokenMessage: null,
  needsToken: false,
};

function repositoryDetailReducer(
  state: RepositoryDetailState,
  action: RepositoryDetailAction
): RepositoryDetailState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, error: null, fileError: null };
    case "FETCH_REPOSITORY_SUCCESS":
      return { ...state, repository: action.repository, isLoading: false };
    case "FETCH_FILES_SUCCESS":
      return {
        ...state,
        files: action.files,
        isLoading: false,
        fileError: null,
        needsToken: false,
      };
    case "FETCH_ERROR":
      return { ...state, error: action.error, isLoading: false };
    case "FETCH_FILE_ERROR":
      return {
        ...state,
        fileError: action.error,
        isLoading: false,
        needsToken: action.needsToken || false,
      };
    case "SET_ACCESS_TOKEN":
      return { ...state, accessToken: action.token };
    case "UPDATE_TOKEN_START":
      return { ...state, isUpdatingToken: true, tokenMessage: null };
    case "UPDATE_TOKEN_SUCCESS":
      return {
        ...state,
        isUpdatingToken: false,
        tokenMessage: { type: "success", text: action.message },
        accessToken: "",
        needsToken: false,
      };
    case "UPDATE_TOKEN_ERROR":
      return {
        ...state,
        isUpdatingToken: false,
        tokenMessage: { type: "error", text: action.error },
      };
    case "CLEAR_TOKEN_MESSAGE":
      return { ...state, tokenMessage: null };
    default:
      return state;
  }
}

export function useRepositoryDetail(repoId: string | undefined) {
  const [state, dispatch] = useReducer(repositoryDetailReducer, initialState);
  const queryService = useRepositoryQueryService();
  const commandService = useRepositoryCommandService();

  const fetchRepository = useCallback(async () => {
    if (!repoId) return;

    dispatch({ type: "FETCH_START" });

    try {
      const repository = await queryService.getById(repoId);
      dispatch({ type: "FETCH_REPOSITORY_SUCCESS", repository });
    } catch (err) {
      dispatch({
        type: "FETCH_ERROR",
        error: "Failed to load repository details. Please try again later.",
      });
      console.error("Error fetching repository:", err);
    }
  }, [repoId, queryService]);

  const fetchFiles = useCallback(async () => {
    if (!repoId) return;

    dispatch({ type: "FETCH_START" });

    try {
      const files = await queryService.listFiles(repoId);
      dispatch({ type: "FETCH_FILES_SUCCESS", files });
    } catch (err) {
      // Handle Axios errors
      if (err instanceof AxiosError) {
        if (err.response?.status === 401) {
          dispatch({
            type: "FETCH_FILE_ERROR",
            error: "Authentication required. Please log in.",
            needsToken: false,
          });
          return;
        }

        if (err.response?.status === 400) {
          const errorData = err.response?.data;
          if (errorData?.code === "ACCESS_TOKEN_REQUIRED") {
            dispatch({
              type: "FETCH_FILE_ERROR",
              error: "Access token is required to list repository files",
              needsToken: true,
            });
            return;
          }
        }
      }

      dispatch({
        type: "FETCH_FILE_ERROR",
        error: "Failed to load repository files. Please try again later.",
      });
      console.error("Error fetching files:", err);
    }
  }, [repoId, queryService]);

  useEffect(() => {
    if (repoId) {
      Promise.all([fetchRepository(), fetchFiles()]);
    }
  }, [repoId, fetchRepository, fetchFiles]);

  const handleTokenSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!state.accessToken.trim()) {
      dispatch({
        type: "UPDATE_TOKEN_ERROR",
        error: "Please enter an access token",
      });
      return;
    }

    if (!repoId) return;

    dispatch({ type: "UPDATE_TOKEN_START" });

    try {
      await commandService.updateAccessToken(repoId, state.accessToken);

      dispatch({
        type: "UPDATE_TOKEN_SUCCESS",
        message: "Access token updated successfully!",
      });

      // Refetch files
      fetchFiles();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      dispatch({ type: "UPDATE_TOKEN_ERROR", error: message });
    }
  };

  const setAccessToken = (token: string) => {
    dispatch({ type: "SET_ACCESS_TOKEN", token });
  };

  const markdownFiles = state.files.filter(
    (file) => file.type === "file" && file.path.toLowerCase().endsWith(".md")
  );

  return {
    ...state,
    handleTokenSubmit,
    setAccessToken,
    markdownFiles,
  };
}
