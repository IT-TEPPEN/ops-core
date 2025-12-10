import React, { useState, useEffect } from "react";
import type { User } from "../../types/domain";
import { listUsers } from "../../utils/userApi";

interface GroupMemberSelectorProps {
  existingMemberIds: string[];
  onAdd: (userId: string) => Promise<void>;
}

/**
 * GroupMemberSelector component for adding members to a group
 */
export const GroupMemberSelector: React.FC<GroupMemberSelectorProps> = ({
  existingMemberIds,
  onAdd,
}) => {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isAdding, setIsAdding] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const allUsers = await listUsers();
        setUsers(allUsers);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load users");
      } finally {
        setIsLoading(false);
      }
    };

    fetchUsers();
  }, []);

  const availableUsers = users.filter(
    (user) => !existingMemberIds.includes(user.id)
  );

  const handleAdd = async () => {
    if (!selectedUserId) return;

    setIsAdding(true);
    setError(null);
    try {
      await onAdd(selectedUserId);
      setSelectedUserId("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add member");
    } finally {
      setIsAdding(false);
    }
  };

  if (isLoading) {
    return (
      <div className="text-center text-gray-500 py-4">
        Loading users...
      </div>
    );
  }

  if (availableUsers.length === 0) {
    return (
      <div className="text-center text-gray-500 py-4">
        All users are already members of this group.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <div className="flex gap-3">
        <select
          value={selectedUserId}
          onChange={(e) => setSelectedUserId(e.target.value)}
          className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          disabled={isAdding}
        >
          <option value="">Select a user to add...</option>
          {availableUsers.map((user) => (
            <option key={user.id} value={user.id}>
              {user.name} ({user.email})
            </option>
          ))}
        </select>
        <button
          onClick={handleAdd}
          disabled={!selectedUserId || isAdding}
          className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {isAdding ? "Adding..." : "Add Member"}
        </button>
      </div>
    </div>
  );
};

export default GroupMemberSelector;
