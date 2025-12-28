import { useState, useEffect } from "react";
import { ExecutionRecord } from "@/shared/types/domain";
import {
  getExecutionRecord,
  createExecutionRecord,
  updateExecutionRecordTitle,
  updateExecutionRecordNotes,
  addExecutionStep,
  updateStepNotes,
  completeExecutionRecord,
  failExecutionRecord,
} from "@/shared/api";

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
  const [executionRecord, setExecutionRecord] =
    useState<ExecutionRecord | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Fetch execution record if recordId is provided
  useEffect(() => {
    if (!recordId) return;

    const fetchRecord = async () => {
      try {
        const record = await getExecutionRecord(recordId);
        setExecutionRecord(record);
      } catch (err) {
        console.error("Error fetching execution record:", err);
        setError("Failed to load execution record");
      }
    };

    fetchRecord();
  }, [recordId]);

  const handleStartExecution = async () => {
    if (!docId || !documentVersionId) return;

    setIsCreating(true);
    try {
      const variableValuesList = Object.entries(variableValues).map(
        ([name, value]) => ({ name, value })
      );

      const record = await createExecutionRecord({
        document_id: docId,
        document_version_id: documentVersionId,
        title: executionTitle,
        variable_values: variableValuesList,
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
      const updated = await updateExecutionRecordTitle(
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
      const updated = await updateExecutionRecordNotes(
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

    const updated = await addExecutionStep(
      executionRecord.id,
      stepNumber,
      description
    );
    setExecutionRecord(updated);
  };

  const handleUpdateStepNotes = async (stepNumber: number, notes: string) => {
    if (!executionRecord) return;

    const updated = await updateStepNotes(
      executionRecord.id,
      stepNumber,
      notes
    );
    setExecutionRecord(updated);
  };

  const handleComplete = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await completeExecutionRecord(executionRecord.id);
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
      const updated = await failExecutionRecord(executionRecord.id);
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
