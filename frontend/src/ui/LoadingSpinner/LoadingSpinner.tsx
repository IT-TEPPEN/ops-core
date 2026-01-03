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
    <div className="h-full w-full flex items-center justify-center gap-4">
      <SpinnerIcon
        className="text-blue-600 dark:text-blue-400"
        size={sizeMap[size]}
      />
      {message && <p>{message}</p>}
    </div>
  );
}
