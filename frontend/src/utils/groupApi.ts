/**
 * Group API utilities
 */

import { get, post, put, del, apiRequest } from "./api";
import type { Group } from "../types/domain";

/** Create group request */
export interface CreateGroupRequest {
  name: string;
  description: string;
}

/** Update group request */
export interface UpdateGroupRequest {
  name: string;
  description: string;
}

/** Add member request */
export interface AddMemberRequest {
  user_id: string;
}

/** Remove member request */
export interface RemoveMemberRequest {
  user_id: string;
}

/** List groups response */
export interface ListGroupsResponse {
  groups: Group[];
}

/**
 * Create a new group
 */
export async function createGroup(
  req: CreateGroupRequest,
  signal?: AbortSignal
): Promise<Group> {
  return post<Group>("/groups", req, signal);
}

/**
 * Get group by ID
 */
export async function getGroup(
  groupId: string,
  signal?: AbortSignal
): Promise<Group> {
  return get<Group>(`/groups/${groupId}`, signal);
}

/**
 * List all groups
 */
export async function listGroups(signal?: AbortSignal): Promise<Group[]> {
  const response = await get<ListGroupsResponse>("/groups", signal);
  return response.groups;
}

/**
 * Update group
 */
export async function updateGroup(
  groupId: string,
  req: UpdateGroupRequest,
  signal?: AbortSignal
): Promise<Group> {
  return put<Group>(`/groups/${groupId}`, req, signal);
}

/**
 * Delete group
 */
export async function deleteGroup(
  groupId: string,
  signal?: AbortSignal
): Promise<{ message: string; groupId: string }> {
  return del<{ message: string; groupId: string }>(`/groups/${groupId}`, signal);
}

/**
 * Add member to group
 */
export async function addMember(
  groupId: string,
  req: AddMemberRequest,
  signal?: AbortSignal
): Promise<Group> {
  return post<Group>(`/groups/${groupId}/members`, req, signal);
}

/**
 * Remove member from group
 */
export async function removeMember(
  groupId: string,
  req: RemoveMemberRequest,
  signal?: AbortSignal
): Promise<Group> {
  // Note: DELETE with body - using apiRequest directly
  const response = await apiRequest<Group>(`/groups/${groupId}/members`, {
    method: "DELETE",
    body: req,
    signal,
  });
  return response.data;
}

/**
 * Get groups for a user
 */
export async function getUserGroups(
  userId: string,
  signal?: AbortSignal
): Promise<Group[]> {
  const response = await get<ListGroupsResponse>(`/users/${userId}/groups`, signal);
  return response.groups;
}
