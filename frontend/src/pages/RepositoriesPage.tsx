import { useReducer } from "react";
import { RepositoryList } from "../features/repository";
import { UI_Form_Input, UI_Form_Submit } from "../ui/form";
import { UI_Form_Field } from "../components";

interface RepositoryFormState {
  newRepoUrl: string;
  isSubmitting: boolean;
  submitMessage: {
    type: "success" | "error";
    text: string;
  } | null;
}

type RepositoryFormAction =
  | { type: "SET_REPO_URL"; payload: string }
  | { type: "SET_SUBMITTING"; payload: boolean }
  | {
      type: "SET_SUBMIT_MESSAGE";
      payload: { type: "success" | "error"; text: string } | null;
    }
  | { type: "RESET_FORM" };

function repositoryFormReducer(
  state: RepositoryFormState,
  action: RepositoryFormAction
): RepositoryFormState {
  switch (action.type) {
    case "SET_REPO_URL":
      return { ...state, newRepoUrl: action.payload };
    case "SET_SUBMITTING":
      return { ...state, isSubmitting: action.payload };
    case "SET_SUBMIT_MESSAGE":
      return { ...state, submitMessage: action.payload };
    case "RESET_FORM":
      return { ...state, newRepoUrl: "", submitMessage: null };
    default:
      return state;
  }
}

function RepositoriesPage() {
  const [state, dispatch] = useReducer(repositoryFormReducer, {
    newRepoUrl: "",
    isSubmitting: false,
    submitMessage: null,
  });

  // API base URL - directly use the base URL to avoid recalculation
  const apiHost = import.meta.env.VITE_API_HOST || window.location.host;
  const apiUrl = `${window.location.protocol}//${apiHost}/api/v1`;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    dispatch({ type: "SET_SUBMITTING", payload: true });
    dispatch({ type: "SET_SUBMIT_MESSAGE", payload: null });

    try {
      const response = await fetch(`${apiUrl}/repositories`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: state.newRepoUrl }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to register repository");
      }

      dispatch({
        type: "SET_SUBMIT_MESSAGE",
        payload: {
          type: "success",
          text: "Repository registered successfully!",
        },
      });
      dispatch({ type: "RESET_FORM" });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unknown error occurred";
      dispatch({
        type: "SET_SUBMIT_MESSAGE",
        payload: { type: "error", text: message },
      });
    } finally {
      dispatch({ type: "SET_SUBMITTING", payload: false });
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Repository Management</h1>

      {/* Registration Form */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Register New Repository</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <UI_Form_Field label="Repository URL" name="repoUrl" required>
            <UI_Form_Input
              id="repoUrl"
              type="text"
              placeholder="https://github.com/username/repo.git"
              value={state.newRepoUrl}
              onChange={(e) =>
                dispatch({ type: "SET_REPO_URL", payload: e.target.value })
              }
              required
            />
          </UI_Form_Field>
          <UI_Form_Submit
            isSubmitting={state.isSubmitting}
            label="Register Repository"
          />
        </form>

        {state.submitMessage && (
          <div
            className={`mt-4 p-3 rounded ${
              state.submitMessage.type === "success"
                ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
            }`}
          >
            {state.submitMessage.text}
          </div>
        )}
      </div>

      <RepositoryList />
    </div>
  );
}

export default RepositoriesPage;
