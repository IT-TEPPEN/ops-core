interface AlertProps {
  type: "success" | "error" | "warning" | "info";
  children: React.ReactNode;
}

export function Alert({ type, children }: AlertProps) {
  const styles = {
    success:
      "bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200",
    error: "bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200",
    warning:
      "bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200",
    info: "bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200",
  };

  return <div className={`p-4 rounded-lg ${styles[type]}`}>{children}</div>;
}
