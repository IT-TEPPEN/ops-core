import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import type { Group, User } from "../types/domain";
import {
  getGroup,
  updateGroup,
  deleteGroup,
  addMember,
  removeMember,
} from "../api/groupApi";
import { listUsers } from "../api/userApi";
import GroupForm from "../components/Form/GroupForm";
import GroupMemberList from "../components/Display/GroupMemberList";
import GroupMemberSelector from "../components/Form/GroupMemberSelector";
import { formatDateLocale } from "../utils/date";

/**
 * GroupDetailPage displays group details and allows editing
 */
const GroupDetailPage: React.FC = () => {
  const { groupId } = useParams<{ groupId: string }>();
  const navigate = useNavigate();
  const [group, setGroup] = useState<Group | null>(null);
  const [members, setMembers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showMemberSelector, setShowMemberSelector] = useState(false);

  useEffect(() => {
    if (groupId) {
      fetchGroup();
    }
  }, [groupId]);

  const fetchGroup = async () => {
    if (!groupId) return;

    setIsLoading(true);
    setError(null);
    try {
      const data = await getGroup(groupId);
      setGroup(data);

      // Fetch member details
      const allUsers = await listUsers();
      const groupMembers = allUsers.filter((user) =>
        data.member_ids.includes(user.id)
      );
      setMembers(groupMembers);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load group");
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdate = async (data: { name: string; description: string }) => {
    if (!groupId) return;

    try {
      const updated = await updateGroup(groupId, data);
      setGroup(updated);
      setIsEditing(false);
    } catch (err) {
      throw err; // Let the form handle the error
    }
  };

  const handleDelete = async () => {
    if (!groupId || !window.confirm("Are you sure you want to delete this group?")) {
      return;
    }

    try {
      await deleteGroup(groupId);
      navigate("/groups");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete group");
    }
  };

  const handleAddMember = async (userId: string) => {
    if (!groupId) return;

    try {
      const updated = await addMember(groupId, { user_id: userId });
      setGroup(updated);
      setShowMemberSelector(false);
      
      // Update member list with newly added user
      const allUsers = await listUsers();
      const groupMembers = allUsers.filter((user) =>
        updated.member_ids.includes(user.id)
      );
      setMembers(groupMembers);
    } catch (err) {
      throw err; // Let the selector handle the error
    }
  };

  const handleRemoveMember = async (userId: string) => {
    if (!groupId || !window.confirm("Are you sure you want to remove this member?")) {
      return;
    }

    try {
      const updated = await removeMember(groupId, { user_id: userId });
      setGroup(updated);
      
      // Update member list by filtering out removed user
      setMembers(members.filter((member) => member.id !== userId));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to remove member");
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto p-6 text-center text-gray-500">
        Loading group details...
      </div>
    );
  }

  if (error && !group) {
    return (
      <div className="max-w-5xl mx-auto p-6">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 px-4 py-3 rounded">
          {error}
          <button
            onClick={() => navigate("/groups")}
            className="ml-4 text-sm underline hover:no-underline"
          >
            Back to Groups
          </button>
        </div>
      </div>
    );
  }

  if (!group) {
    return null;
  }

  return (
    <div className="max-w-5xl mx-auto p-6">
      <div className="mb-6">
        <button
          onClick={() => navigate("/groups")}
          className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 mb-4"
        >
          ← Back to Groups
        </button>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {isEditing ? (
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
              Edit Group
            </h2>
            <GroupForm
              initialData={group}
              onSubmit={handleUpdate}
              onCancel={() => setIsEditing(false)}
              submitLabel="Update Group"
            />
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                  {group.name}
                </h1>
                {group.description && (
                  <p className="text-gray-600 dark:text-gray-400 mb-3">
                    {group.description}
                  </p>
                )}
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  <p>Created: {formatDateLocale(group.created_at)}</p>
                  <p>Updated: {formatDateLocale(group.updated_at)}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setIsEditing(true)}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
                >
                  Edit
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            Members ({members.length})
          </h2>
          {!showMemberSelector && (
            <button
              onClick={() => setShowMemberSelector(true)}
              className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition"
            >
              Add Member
            </button>
          )}
        </div>

        {showMemberSelector && (
          <div className="mb-4">
            <GroupMemberSelector
              existingMemberIds={group.member_ids}
              onAdd={handleAddMember}
            />
            <button
              onClick={() => setShowMemberSelector(false)}
              className="mt-2 text-sm text-gray-600 dark:text-gray-400 hover:underline"
            >
              Cancel
            </button>
          </div>
        )}

        <GroupMemberList members={members} onRemove={handleRemoveMember} />
      </div>
    </div>
  );
};

export default GroupDetailPage;
