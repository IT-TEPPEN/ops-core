/**
 * Group API Client
 */

import { V1ApiClient } from "./client";
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
interface ListGroupsResponse {
  groups: Group[];
}

/** Delete group response */
interface DeleteGroupResponse {
  message: string;
  groupId: string;
}

export class GroupApi extends V1ApiClient {
  constructor() {
    super("/groups");
  }

  /**
   * Create a new group
   */
  async createGroup(req: CreateGroupRequest): Promise<Group> {
    return this.post<Group, CreateGroupRequest>("", req);
  }

  /**
   * Get group by ID
   */
  async getGroup(groupId: string): Promise<Group> {
    return this.get<Group>(`/${groupId}`);
  }

  /**
   * List all groups
   */
  async listGroups(): Promise<Group[]> {
    const response = await this.get<ListGroupsResponse>("");
    return response.groups;
  }

  /**
   * Update group
   */
  async updateGroup(groupId: string, req: UpdateGroupRequest): Promise<Group> {
    return this.put<Group, UpdateGroupRequest>(`/${groupId}`, req);
  }

  /**
   * Delete group
   */
  async deleteGroup(groupId: string): Promise<DeleteGroupResponse> {
    return this.delete<DeleteGroupResponse>(`/${groupId}`);
  }

  /**
   * Add member to group
   */
  async addMember(groupId: string, req: AddMemberRequest): Promise<Group> {
    return this.post<Group, AddMemberRequest>(`/${groupId}/members`, req);
  }

  /**
   * Remove member from group
   */
  async removeMember(
    groupId: string,
    req: RemoveMemberRequest
  ): Promise<Group> {
    // Note: DELETE with body - using custom implementation
    return this.deleteWithBody<Group, RemoveMemberRequest>(
      `/${groupId}/members`,
      req
    );
  }

  /**
   * Get groups for a user (different base path)
   */
  async getUserGroups(userId: string): Promise<Group[]> {
    // This uses a different endpoint, need to use a different approach
    const userGroupApi = new V1ApiClient("/users");
    // Access protected get method via type assertion
    const response = await (
      userGroupApi as unknown as { get<T>(endpoint: string): Promise<T> }
    ).get<ListGroupsResponse>(`/${userId}/groups`);
    return response.groups;
  }
}

// Export singleton instance for convenience
export const groupApi = new GroupApi();
