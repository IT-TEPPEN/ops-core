import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import type { Group } from "../../types/domain";
import { getUserGroups } from "../../api";

interface UserGroupBadgeProps {
  userId: string;
  maxDisplay?: number;
}

/**
 * UserGroupBadge component displays badges for user's groups
 */
export const UserGroupBadge: React.FC<UserGroupBadgeProps> = ({
  userId,
  maxDisplay = 3,
}) => {
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const userGroups = await getUserGroups(userId);
        setGroups(userGroups);
      } catch (err) {
        console.error("Failed to load user groups:", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchGroups();
  }, [userId]);

  if (isLoading) {
    return (
      <span className="text-sm text-gray-500 dark:text-gray-400">
        Loading groups...
      </span>
    );
  }

  if (groups.length === 0) {
    return (
      <span className="text-sm text-gray-500 dark:text-gray-400">
        No groups
      </span>
    );
  }

  const displayGroups = groups.slice(0, maxDisplay);
  const remainingCount = groups.length - maxDisplay;

  return (
    <div className="flex flex-wrap gap-2">
      {displayGroups.map((group) => (
        <Link
          key={group.id}
          to={`/groups/${group.id}`}
          className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 hover:bg-blue-200 dark:hover:bg-blue-800 transition"
        >
          {group.name}
        </Link>
      ))}
      {remainingCount > 0 && (
        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
          +{remainingCount} more
        </span>
      )}
    </div>
  );
};

export default UserGroupBadge;
