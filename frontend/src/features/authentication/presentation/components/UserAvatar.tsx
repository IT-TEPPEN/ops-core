interface UserAvatarProps {
  name: string;
  pictureUrl?: string;
  size?: "sm" | "md" | "lg";
}

const sizeClasses = {
  sm: "w-8 h-8 text-xs",
  md: "w-10 h-10 text-sm",
  lg: "w-12 h-12 text-base",
};

export function UserAvatar({ name, pictureUrl, size = "md" }: UserAvatarProps) {
  const sizeClass = sizeClasses[size];

  // Get initials from name (first letter of first and last name)
  const getInitials = (name: string): string => {
    const parts = name.trim().split(/\s+/);
    if (parts.length === 0) return "?";
    if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
    return (
      parts[0].charAt(0).toUpperCase() +
      parts[parts.length - 1].charAt(0).toUpperCase()
    );
  };

  if (pictureUrl) {
    return (
      <img
        src={pictureUrl}
        alt={name}
        className={`${sizeClass} rounded-full object-cover border-2 border-gray-200 dark:border-gray-700`}
      />
    );
  }

  // Fallback to initials
  return (
    <div
      className={`${sizeClass} rounded-full bg-blue-600 dark:bg-blue-500 text-white font-medium flex items-center justify-center border-2 border-gray-200 dark:border-gray-700`}
    >
      {getInitials(name)}
    </div>
  );
}
