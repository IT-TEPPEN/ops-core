import { useState, useEffect } from "react";
import type { ExecutionRecordViewData } from "../application/dto";
import {
  useExecutionQueryService,
  useExecutionCommandService,
} from "../presentation/contexts";

interface UseExecutionRecordProps {
  docId: string | undefined;
  recordId: string | undefined;
  documentVersionId: string | undefined;
  variableValues: Record<string, string | number | boolean>;
  executionTitle: string;
}

export function useExecutionRecord({
  docId,
  recordId,
  documentVersionId,
  variableValues,
  executionTitle,
}: UseExecutionRecordProps) {
  const executionQueryService = useExecutionQueryService();
  const executionCommandService = useExecutionCommandService();

  const [executionRecord, setExecutionRecord] =
    useState<ExecutionRecordViewData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Fetch execution record if recordId is provided
  useEffect(() => {
    if (!recordId) return;

    const fetchRecord = async () => {
      try {
        const record = await executionQueryService.getById(recordId);
        setExecutionRecord(record);
      } catch (err) {
        console.error("Error fetching execution record:", err);
        setError("Failed to load execution record");
      }
    };

    fetchRecord();
  }, [recordId, executionQueryService]);

  const handleStartExecution = async () => {
    if (!docId || !documentVersionId) return;

    setIsCreating(true);
    try {
      const variableValuesList = Object.entries(variableValues).map(
        ([name, value]) => ({ name, value })
      );

      const record = await executionCommandService.create({
        documentId: docId,
        documentVersionId: documentVersionId,
        title: executionTitle,
        variableValues: variableValuesList,
      });

      setExecutionRecord(record);
      setError(null);
    } catch (err) {
      console.error("Failed to create execution record:", err);
      setError("Failed to start execution");
    } finally {
      setIsCreating(false);
    }
  };

  const handleUpdateTitle = async (newTitle: string) => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await executionCommandService.updateTitle(
        executionRecord.id,
        newTitle
      );
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to update title:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleUpdateNotes = async (notes: string) => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await executionCommandService.updateNotes(
        executionRecord.id,
        notes
      );
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to update notes:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleAddStep = async (stepNumber: number, description: string) => {
    if (!executionRecord) return;

    const updated = await executionCommandService.addStep(executionRecord.id, {
      stepNumber,
      description,
    });
    setExecutionRecord(updated);
  };

  const handleUpdateStepNotes = async (stepNumber: number, notes: string) => {
    if (!executionRecord) return;

    const updated = await executionCommandService.updateStepNotes(
      executionRecord.id,
      {
        stepNumber,
        notes,
      }
    );
    setExecutionRecord(updated);
  };

  const handleComplete = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await executionCommandService.complete(executionRecord.id);
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to complete execution:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleFail = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await executionCommandService.fail(executionRecord.id);
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to fail execution:", err);
    } finally {
      setIsSaving(false);
    }
  };

  return {
    executionRecord,
    error,
    isCreating,
    isSaving,
    handleStartExecution,
    handleUpdateTitle,
    handleUpdateNotes,
    handleAddStep,
    handleUpdateStepNotes,
    handleComplete,
    handleFail,
  };
}
