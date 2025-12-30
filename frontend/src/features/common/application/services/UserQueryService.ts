import type { User, Group } from "@/shared/types/domain";

/**
 * User Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface UserQueryService {
  /**
   * List all users.
   */
  listUsers(): Promise<User[]>;

  /**
   * Get groups for a specific user.
   */
  getUserGroups(userId: string): Promise<Group[]>;
}
