import React from "react";
import { Link } from "react-router-dom";
import type { Group } from "../../types/domain";
import { formatDateLocale } from "../../utils/date";

interface GroupListProps {
  groups: Group[];
}

/**
 * GroupList component displays a list of groups
 */
export const GroupList: React.FC<GroupListProps> = ({ groups }) => {
  if (groups.length === 0) {
    return (
      <div className="text-center text-gray-500 py-8">
        No groups found. Create your first group to get started.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {groups.map((group) => (
        <Link
          key={group.id}
          to={`/groups/${group.id}`}
          className="block bg-white dark:bg-gray-800 p-6 rounded-lg shadow hover:shadow-lg transition-shadow border border-gray-200 dark:border-gray-700"
        >
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
                {group.name}
              </h3>
              {group.description && (
                <p className="text-gray-600 dark:text-gray-400 mb-3">
                  {group.description}
                </p>
              )}
              <div className="flex gap-4 text-sm text-gray-500 dark:text-gray-400">
                <span>Members: {group.member_ids.length}</span>
                <span>Created: {formatDateLocale(group.created_at)}</span>
              </div>
            </div>
            <div className="ml-4">
              <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                {group.member_ids.length} {group.member_ids.length === 1 ? "member" : "members"}
              </span>
            </div>
          </div>
        </Link>
      ))}
    </div>
  );
};

export default GroupList;
