import type { Group } from "@/shared/types/domain";

/**
 * Request for creating a new group.
 */
export interface CreateGroupRequest {
  name: string;
  description: string;
}

/**
 * Request for updating a group.
 */
export interface UpdateGroupRequest {
  name: string;
  description: string;
}

/**
 * Request for adding/removing a member.
 */
export interface MemberRequest {
  user_id: string;
}

/**
 * Group Command Service interface for write operations (POST/PUT/DELETE).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface GroupCommandService {
  /**
   * Create a new group.
   */
  create(request: CreateGroupRequest): Promise<Group>;

  /**
   * Update an existing group.
   */
  update(groupId: string, request: UpdateGroupRequest): Promise<Group>;

  /**
   * Delete a group.
   */
  deleteGroup(groupId: string): Promise<void>;

  /**
   * Add a member to a group.
   */
  addMember(groupId: string, request: MemberRequest): Promise<Group>;

  /**
   * Remove a member from a group.
   */
  removeMember(groupId: string, request: MemberRequest): Promise<Group>;
}
