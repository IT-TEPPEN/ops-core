import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useGroupDetail } from "@/features/group/hooks/useGroupDetail";
import { GroupInfoCard } from "@/features/group/components/GroupInfoCard";
import { MembersPanel } from "@/features/group/components/MembersPanel";
import { LoadingSpinner } from "@/ui";

const GroupDetailPage: React.FC = () => {
  const { groupId } = useParams<{ groupId: string }>();
  const navigate = useNavigate();

  const {
    group,
    members,
    isLoading,
    isEditing,
    showMemberSelector,
    error,
    handleUpdate,
    handleDelete,
    handleAddMember,
    handleRemoveMember,
    setEditing,
    setShowMemberSelector,
  } = useGroupDetail(groupId);

  if (isLoading) {
    return <LoadingSpinner message="Loading..." />;
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

        <GroupInfoCard
          group={group}
          isEditing={isEditing}
          onUpdate={handleUpdate}
          onDelete={handleDelete}
          onEditToggle={setEditing}
        />
      </div>

      <MembersPanel
        members={members}
        existingMemberIds={group.member_ids}
        showMemberSelector={showMemberSelector}
        onAddMember={handleAddMember}
        onRemoveMember={handleRemoveMember}
        onToggleSelector={setShowMemberSelector}
      />
    </div>
  );
};

export default GroupDetailPage;
