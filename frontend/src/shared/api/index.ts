/**
 * API module exports
 */

// Export API client utilities
export { V1ApiClient, ApiRequestError } from "./client";

// Export API client classes and instances
export * from "./authApi";
export * from "./repositoryApi";
export * from "./userApiClient";
export * from "./groupApiClient";
export * from "./executionRecordApiClient";
export * from "./gitProviderApi";

// Re-export instance methods as standalone functions for backwards compatibility
import { userApi } from "./userApiClient";
import { groupApi } from "./groupApiClient";
import { executionRecordApi } from "./executionRecordApiClient";

// User API functions
export const createUser = userApi.createUser.bind(userApi);
export const getUser = userApi.getUser.bind(userApi);
export const listUsers = userApi.listUsers.bind(userApi);
export const updateUser = userApi.updateUser.bind(userApi);
export const deleteUser = userApi.deleteUser.bind(userApi);

// Group API functions
export const createGroup = groupApi.createGroup.bind(groupApi);
export const getGroup = groupApi.getGroup.bind(groupApi);
export const listGroups = groupApi.listGroups.bind(groupApi);
export const updateGroup = groupApi.updateGroup.bind(groupApi);
export const deleteGroup = groupApi.deleteGroup.bind(groupApi);
export const addMember = groupApi.addMember.bind(groupApi);
export const removeMember = groupApi.removeMember.bind(groupApi);
export const getUserGroups = groupApi.getUserGroups.bind(groupApi);

// Execution Record API functions
export const createExecutionRecord =
  executionRecordApi.createExecutionRecord.bind(executionRecordApi);
export const getExecutionRecord =
  executionRecordApi.getExecutionRecord.bind(executionRecordApi);
export const searchExecutionRecords =
  executionRecordApi.searchExecutionRecords.bind(executionRecordApi);
export const updateExecutionRecordTitle =
  executionRecordApi.updateExecutionRecordTitle.bind(executionRecordApi);
export const updateExecutionRecordNotes =
  executionRecordApi.updateExecutionRecordNotes.bind(executionRecordApi);
export const addExecutionStep =
  executionRecordApi.addExecutionStep.bind(executionRecordApi);
export const updateStepNotes =
  executionRecordApi.updateStepNotes.bind(executionRecordApi);
export const completeExecutionRecord =
  executionRecordApi.completeExecutionRecord.bind(executionRecordApi);
export const failExecutionRecord =
  executionRecordApi.failExecutionRecord.bind(executionRecordApi);
export const updateExecutionRecordAccessScope =
  executionRecordApi.updateExecutionRecordAccessScope.bind(executionRecordApi);
export const deleteExecutionRecord =
  executionRecordApi.deleteExecutionRecord.bind(executionRecordApi);
