import { useState } from "react";
import { ExecutionStep } from "../../types/domain";

export interface ExecutionStepPanelProps {
  steps: ExecutionStep[];
  onAddStep: (stepNumber: number, description: string) => Promise<void>;
  onUpdateStepNotes: (stepNumber: number, notes: string) => Promise<void>;
}

export function ExecutionStepPanel({
  steps,
  onAddStep,
  onUpdateStepNotes,
}: ExecutionStepPanelProps) {
  const [newStepNumber, setNewStepNumber] = useState<number>(steps.length + 1);
  const [newStepDescription, setNewStepDescription] = useState("");
  const [editingStep, setEditingStep] = useState<number | null>(null);
  const [editNotes, setEditNotes] = useState("");
  const [isAdding, setIsAdding] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  const handleAddStep = async () => {
    if (!newStepDescription.trim()) return;

    setIsAdding(true);
    try {
      await onAddStep(newStepNumber, newStepDescription);
      setNewStepNumber(newStepNumber + 1);
      setNewStepDescription("");
    } catch (error) {
      console.error("Failed to add step:", error);
    } finally {
      setIsAdding(false);
    }
  };

  const handleUpdateNotes = async (stepNumber: number) => {
    setIsUpdating(true);
    try {
      await onUpdateStepNotes(stepNumber, editNotes);
      setEditingStep(null);
      setEditNotes("");
    } catch (error) {
      console.error("Failed to update step notes:", error);
    } finally {
      setIsUpdating(false);
    }
  };

  const startEditing = (step: ExecutionStep) => {
    setEditingStep(step.step_number);
    setEditNotes(step.notes || "");
  };

  const cancelEditing = () => {
    setEditingStep(null);
    setEditNotes("");
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold">Execution Steps</h3>

      {/* Step List */}
      <div className="space-y-3">
        {steps.length === 0 ? (
          <p className="text-sm text-gray-500 dark:text-gray-400">
            No steps added yet. Add a step to track your progress.
          </p>
        ) : (
          steps.map((step) => (
            <div
              key={step.id}
              className="p-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg"
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <span className="inline-block px-2 py-1 text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded">
                    Step {step.step_number}
                  </span>
                  <p className="mt-1 text-sm font-medium">{step.description}</p>
                </div>
              </div>

              {editingStep === step.step_number ? (
                <div className="mt-2 space-y-2">
                  <textarea
                    value={editNotes}
                    onChange={(e) => setEditNotes(e.target.value)}
                    placeholder="Add notes for this step..."
                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
                    rows={3}
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleUpdateNotes(step.step_number)}
                      disabled={isUpdating}
                      className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                    >
                      {isUpdating ? "Saving..." : "Save Notes"}
                    </button>
                    <button
                      onClick={cancelEditing}
                      className="px-3 py-1 text-sm bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 rounded hover:bg-gray-400 dark:hover:bg-gray-500"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              ) : (
                <div className="mt-2">
                  {step.notes ? (
                    <p className="text-sm text-gray-600 dark:text-gray-300 whitespace-pre-wrap">
                      {step.notes}
                    </p>
                  ) : (
                    <p className="text-sm text-gray-400 dark:text-gray-500 italic">
                      No notes
                    </p>
                  )}
                  <button
                    onClick={() => startEditing(step)}
                    className="mt-2 text-sm text-blue-600 dark:text-blue-400 hover:underline"
                  >
                    {step.notes ? "Edit Notes" : "Add Notes"}
                  </button>
                </div>
              )}

              <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
                Executed at: {new Date(step.executed_at).toLocaleString()}
              </p>
            </div>
          ))
        )}
      </div>

      {/* Add New Step */}
      <div className="p-3 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
        <h4 className="text-sm font-semibold mb-2">Add New Step</h4>
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <label className="text-sm font-medium whitespace-nowrap">
              Step #
            </label>
            <input
              type="number"
              min="1"
              value={newStepNumber}
              onChange={(e) => setNewStepNumber(parseInt(e.target.value) || 1)}
              className="w-20 px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <input
            type="text"
            value={newStepDescription}
            onChange={(e) => setNewStepDescription(e.target.value)}
            placeholder="Step description..."
            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
          />
          <button
            onClick={handleAddStep}
            disabled={isAdding || !newStepDescription.trim()}
            className="w-full px-4 py-2 text-sm bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isAdding ? "Adding..." : "Add Step"}
          </button>
        </div>
      </div>
    </div>
  );
}
