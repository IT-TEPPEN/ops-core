import React from "react";
import type { User } from "../../types/domain";

interface GroupMemberListProps {
  members: User[];
  onRemove?: (userId: string) => void;
  isLoading?: boolean;
}

/**
 * GroupMemberList component displays a list of group members
 */
export const GroupMemberList: React.FC<GroupMemberListProps> = ({
  members,
  onRemove,
  isLoading = false,
}) => {
  if (isLoading) {
    return (
      <div className="text-center text-gray-500 py-8">
        Loading members...
      </div>
    );
  }

  if (members.length === 0) {
    return (
      <div className="text-center text-gray-500 py-8">
        No members in this group yet.
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {members.map((member) => (
        <div
          key={member.id}
          className="flex items-center justify-between p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700"
        >
          <div className="flex-1">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold">
                {member.name.charAt(0).toUpperCase()}
              </div>
              <div>
                <h4 className="text-sm font-medium text-gray-900 dark:text-gray-100">
                  {member.name}
                </h4>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {member.email}
                </p>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span className="text-xs px-2 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 capitalize">
              {member.role}
            </span>
            {onRemove && (
              <button
                onClick={() => onRemove(member.id)}
                className="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 text-sm font-medium"
                title="Remove member"
              >
                Remove
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};

export default GroupMemberList;
