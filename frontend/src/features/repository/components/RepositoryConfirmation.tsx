import { GitRepository, GitProvider } from "@/shared/api/gitProviderApi";
import { Card, Button } from "@/ui";
import { useRepositoryRegistration } from "../hooks/useRepositoryRegistration";

interface RepositoryConfirmationProps {
  repository: GitRepository;
  provider: GitProvider;
  onCancel: () => void;
  onSuccess?: () => void;
}

export function RepositoryConfirmation({
  repository,
  provider,
  onCancel,
  onSuccess,
}: RepositoryConfirmationProps) {
  const { isSubmitting, message, handleSubmitRepository } =
    useRepositoryRegistration();

  const handleRegister = async () => {
    const success = await handleSubmitRepository({
      provider,
      url: repository.cloneUrl,
    });
    if (success) {
      onSuccess?.();
    }
  };
  return (
    <Card>
      <h3 className="text-lg font-semibold mb-4">Confirm Registration</h3>
      <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg mb-4">
        <div className="flex items-center gap-3">
          <img
            src={repository.owner.avatarUrl}
            alt={repository.owner.login}
            className="w-10 h-10 rounded-full"
          />
          <div>
            <p className="font-medium">{repository.fullName}</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {repository.cloneUrl}
            </p>
          </div>
        </div>
      </div>

      <div className="flex gap-3">
        <Button
          variant="primary"
          onClick={handleRegister}
          isLoading={isSubmitting}
        >
          {isSubmitting ? "Registering..." : "Register Repository"}
        </Button>
        <Button variant="secondary" onClick={onCancel} disabled={isSubmitting}>
          Cancel
        </Button>
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
    </Card>
  );
}
