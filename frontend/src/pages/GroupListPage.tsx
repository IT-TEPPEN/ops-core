import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import type { Group } from "../types/domain";
import { listGroups, createGroup } from "../utils/groupApi";
import GroupList from "../components/Display/GroupList";
import GroupForm from "../components/Form/GroupForm";

/**
 * GroupListPage displays all groups and allows creating new ones
 */
const GroupListPage: React.FC = () => {
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchGroups();
  }, []);

  const fetchGroups = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await listGroups();
      setGroups(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load groups");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreate = async (data: { name: string; description: string }) => {
    try {
      const newGroup = await createGroup(data);
      setShowCreateForm(false);
      navigate(`/groups/${newGroup.id}`);
    } catch (err) {
      throw err; // Let the form handle the error
    }
  };

  return (
    <div className="max-w-5xl mx-auto p-6">
      <div className="mb-6">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            Groups
          </h1>
          {!showCreateForm && (
            <button
              onClick={() => setShowCreateForm(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
            >
              Create Group
            </button>
          )}
        </div>

        {showCreateForm && (
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow mb-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-gray-100">
              Create New Group
            </h2>
            <GroupForm
              onSubmit={handleCreate}
              onCancel={() => setShowCreateForm(false)}
              submitLabel="Create Group"
            />
          </div>
        )}
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 px-4 py-3 rounded mb-4">
          {error}
          <button
            onClick={fetchGroups}
            className="ml-4 text-sm underline hover:no-underline"
          >
            Retry
          </button>
        </div>
      )}

      {isLoading ? (
        <div className="text-center text-gray-500 py-8">Loading groups...</div>
      ) : (
        <GroupList groups={groups} />
      )}
    </div>
  );
};

export default GroupListPage;
