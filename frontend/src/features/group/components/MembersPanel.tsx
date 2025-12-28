import { User } from "@/shared/types/domain";
import GroupMemberList from "@/features/common/components/Display/GroupMemberList";
import GroupMemberSelector from "@/features/common/components/Form/GroupMemberSelector";

interface MembersPanelProps {
  members: User[];
  existingMemberIds: string[];
  showMemberSelector: boolean;
  onAddMember: (userId: string) => Promise<void>;
  onRemoveMember: (userId: string) => Promise<void>;
  onToggleSelector: (show: boolean) => void;
}

export function MembersPanel({
  members,
  existingMemberIds,
  showMemberSelector,
  onAddMember,
  onRemoveMember,
  onToggleSelector,
}: MembersPanelProps) {
  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Members ({members.length})
        </h2>
        {!showMemberSelector && (
          <button
            onClick={() => onToggleSelector(true)}
            className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition"
          >
            Add Member
          </button>
        )}
      </div>

      {showMemberSelector && (
        <div className="mb-4">
          <GroupMemberSelector
            existingMemberIds={existingMemberIds}
            onAdd={onAddMember}
          />
          <button
            onClick={() => onToggleSelector(false)}
            className="mt-2 text-sm text-gray-600 dark:text-gray-400 hover:underline"
          >
            Cancel
          </button>
        </div>
      )}

      <GroupMemberList members={members} onRemove={onRemoveMember} />
    </div>
  );
}
