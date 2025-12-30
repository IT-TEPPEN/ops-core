import type { Group, User } from "@/shared/types/domain";

/**
 * Group Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface GroupQueryService {
  /**
   * Get group details by ID.
   */
  getById(groupId: string): Promise<Group>;

  /**
   * List all groups.
   */
  list(): Promise<Group[]>;

  /**
   * List all users (for member management).
   */
  listUsers(): Promise<User[]>;

  /**
   * Get groups for a specific user.
   */
  getUserGroups(userId: string): Promise<Group[]>;
}
