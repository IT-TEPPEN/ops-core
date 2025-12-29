import { SpinnerIcon } from "../Icon";

interface LoadingSpinnerProps {
  message?: string;
  size?: "sm" | "md" | "lg";
}

export function LoadingSpinner({ message, size = "md" }: LoadingSpinnerProps) {
  const sizeMap = {
    sm: 16,
    md: 24,
    lg: 32,
  };

  return (
    <div className="flex items-center justify-center py-4">
      <SpinnerIcon className="text-blue-600 mr-3" size={sizeMap[size]} />
      {message && <span>{message}</span>}
    </div>
  );
}
