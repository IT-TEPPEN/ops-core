import { formatDateLocale } from "@/shared/utils/date";
import { Group } from "@/shared/types/domain";
import GroupForm from "@/features/common/components/Form/GroupForm";

interface GroupInfoCardProps {
  group: Group;
  isEditing: boolean;
  onUpdate: (data: { name: string; description: string }) => Promise<void>;
  onDelete: () => Promise<void>;
  onEditToggle: (isEditing: boolean) => void;
}

export function GroupInfoCard({
  group,
  isEditing,
  onUpdate,
  onDelete,
  onEditToggle,
}: GroupInfoCardProps) {
  if (isEditing) {
    return (
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
          Edit Group
        </h2>
        <GroupForm
          initialData={group}
          onSubmit={onUpdate}
          onCancel={() => onEditToggle(false)}
          submitLabel="Update Group"
        />
      </div>
    );
  }

  return (
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
            onClick={() => onEditToggle(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
          >
            Edit
          </button>
          <button
            onClick={onDelete}
            className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
