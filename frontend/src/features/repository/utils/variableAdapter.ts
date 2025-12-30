import type { DocumentVariable } from "../types/repository";
import type { VariableDefinition } from "@/shared/types/domain";

/**
 * Converts a DocumentVariable to VariableDefinition.
 * Handles the camelCase to snake_case conversion for defaultValue.
 */
export function documentVariableToDefinition(
  docVar: DocumentVariable
): VariableDefinition {
  return {
    name: docVar.name,
    label: docVar.label,
    description: docVar.description ?? undefined,
    type: docVar.type,
    required: docVar.required,
    default_value: docVar.defaultValue ?? undefined,
  };
}

/**
 * Converts an array of DocumentVariables to VariableDefinitions.
 */
export function documentVariablesToDefinitions(
  docVars: DocumentVariable[]
): VariableDefinition[] {
  return docVars.map(documentVariableToDefinition);
}
