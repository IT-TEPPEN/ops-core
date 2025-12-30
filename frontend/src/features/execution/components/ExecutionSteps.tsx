import { ExecutionRecord, ExecutionStep } from "@/shared/types/domain";

interface ExecutionStepsProps {
  executionRecord: ExecutionRecord | null;
  onAddStep: (stepData: { step_number: number; description: string }) => void;
  isSaving: boolean;
}

export function ExecutionSteps({
  executionRecord,
  onAddStep,
  isSaving,
}: ExecutionStepsProps) {
  const handleCheckboxChange = (step: ExecutionStep) => {
    if (!step.executed_at && executionRecord?.status === "in_progress") {
      onAddStep({
        step_number: step.step_number,
        description: step.description,
      });
    }
  };

  if (!executionRecord) {
    return null;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4">Execution Steps</h2>
      <div className="space-y-3">
        {executionRecord.steps.map((step) => (
          <div key={step.id} className="flex items-start gap-3">
            <input
              type="checkbox"
              checked={!!step.executed_at}
              onChange={() => handleCheckboxChange(step)}
              disabled={
                !!step.executed_at ||
                executionRecord.status !== "in_progress" ||
                isSaving
              }
              className="mt-1 w-5 h-5 cursor-pointer disabled:cursor-not-allowed"
            />
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <span className="font-semibold">Step {step.step_number}</span>
                {!!step.executed_at && (
                  <span className="text-xs text-green-600">
                    ✓ Completed at{" "}
                    {new Date(step.executed_at!).toLocaleString()}
                  </span>
                )}
              </div>
              <p className="text-gray-700 mt-1">{step.description}</p>
            </div>
          </div>
        ))}
      </div>
      {executionRecord.steps.length === 0 && (
        <p className="text-gray-500 text-center py-8">
          No steps have been completed yet.
        </p>
      )}
    </div>
  );
}
