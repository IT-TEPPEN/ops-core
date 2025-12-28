import { UseFormRegister, FieldErrors } from "react-hook-form";
import { UI_Form_Field, UI_Form_Input, UI_Form_Submit } from "@/ui";

type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

interface RepositoryFormData {
  provider: GitProvider;
  url: string;
  gitlabUrl?: string;
  gitlabClientId?: string;
  gitlabClientSecret?: string;
}

interface RepositoryRegistrationFormProps {
  register: UseFormRegister<RepositoryFormData>;
  errors: FieldErrors<RepositoryFormData>;
  isSubmitting: boolean;
  onCancel: () => void;
  message?: {
    type: "success" | "error";
    text: string;
  } | null;
}

export function RepositoryRegistrationForm({
  register,
  errors,
  isSubmitting,
  onCancel,
  message,
}: RepositoryRegistrationFormProps) {
  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">
        Step 2: Register Repository
      </h2>
      <div className="space-y-4">
        <UI_Form_Field
          label="Repository URL"
          name="repoUrl"
          error={errors.url?.message}
          required
        >
          <UI_Form_Input
            id="repoUrl"
            type="text"
            placeholder="https://github.com/username/repo.git"
            {...register("url")}
          />
        </UI_Form_Field>

        <div className="flex gap-3">
          <UI_Form_Submit
            isSubmitting={isSubmitting}
            label="Register Repository"
          />
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            disabled={isSubmitting}
          >
            Cancel
          </button>
        </div>
      </div>

      {message && (
        <div
          className={`mt-4 p-3 rounded ${
            message.type === "success"
              ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
              : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
          }`}
        >
          {message.text}
        </div>
      )}
    </div>
  );
}
