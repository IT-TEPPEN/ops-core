import { useReducer } from "react";
import { useNavigate } from "react-router-dom";
import { initiateOAuthFlow, SelfHostedOAuthParams } from "@/shared/utils/oauth";

type RegistrationState = {
  isAuthenticating: boolean;
  isSubmitting: boolean;
  message: {
    type: "success" | "error";
    text: string;
  } | null;
};

type RegistrationAction =
  | { type: "START_AUTH" }
  | { type: "AUTH_ERROR"; error: string }
  | { type: "START_SUBMIT" }
  | { type: "SUBMIT_SUCCESS"; message: string }
  | { type: "SUBMIT_ERROR"; error: string }
  | { type: "CLEAR_MESSAGE" };

const initialState: RegistrationState = {
  isAuthenticating: false,
  isSubmitting: false,
  message: null,
};

function registrationReducer(
  state: RegistrationState,
  action: RegistrationAction
): RegistrationState {
  switch (action.type) {
    case "START_AUTH":
      return { ...state, isAuthenticating: true, message: null };
    case "AUTH_ERROR":
      return {
        ...state,
        isAuthenticating: false,
        message: { type: "error", text: action.error },
      };
    case "START_SUBMIT":
      return { ...state, isSubmitting: true, message: null };
    case "SUBMIT_SUCCESS":
      return {
        ...state,
        isSubmitting: false,
        message: { type: "success", text: action.message },
      };
    case "SUBMIT_ERROR":
      return {
        ...state,
        isSubmitting: false,
        message: { type: "error", text: action.error },
      };
    case "CLEAR_MESSAGE":
      return { ...state, message: null };
    default:
      return state;
  }
}

export function useRepositoryRegistration() {
  const navigate = useNavigate();
  const [state, dispatch] = useReducer(registrationReducer, initialState);

  const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
  const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

  const handleOAuthConnect = async (
    provider: "github" | "gitlab" | "gitlab-self-hosted",
    selfHostedParams?: SelfHostedOAuthParams
  ) => {
    dispatch({ type: "START_AUTH" });

    try {
      await initiateOAuthFlow(provider, selfHostedParams);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to initiate OAuth flow";
      dispatch({ type: "AUTH_ERROR", error: message });
    }
  };

  const handleSubmitRepository = async (data: {
    provider: string;
    url: string;
  }) => {
    dispatch({ type: "START_SUBMIT" });

    try {
      const response = await fetch(`${apiUrl}/repositories`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      const responseData = await response.json();

      if (!response.ok) {
        throw new Error(
          responseData.message || "Failed to register repository"
        );
      }

      dispatch({
        type: "SUBMIT_SUCCESS",
        message: "Repository registered successfully!",
      });

      // 成功後、2秒待ってから一覧ページに遷移
      setTimeout(() => {
        navigate("/repositories");
      }, 2000);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      dispatch({ type: "SUBMIT_ERROR", error: message });
    }
  };

  const clearMessage = () => {
    dispatch({ type: "CLEAR_MESSAGE" });
  };

  return {
    ...state,
    handleOAuthConnect,
    handleSubmitRepository,
    clearMessage,
  };
}
