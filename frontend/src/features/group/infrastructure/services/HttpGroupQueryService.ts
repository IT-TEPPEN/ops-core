import type { GroupQueryService } from "../../application";
import type { Group, User } from "@/shared/types/domain";
import { V1ApiClient } from "@/shared/api/client";

/**
 * HTTP implementation of GroupQueryService.
 * Handles all read operations (GET) for group management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpGroupQueryService
  extends V1ApiClient
  implements GroupQueryService
{
  constructor() {
    super("/groups");
  }

  async getById(groupId: string): Promise<Group> {
    return await this.get<Group>(`/${groupId}`);
  }

  async list(): Promise<Group[]> {
    return await this.get<Group[]>("");
  }

  async listUsers(): Promise<User[]> {
    // Users endpoint is at root level, not under /groups
    return await this.get<User[]>("/users");
  }

  async getUserGroups(userId: string): Promise<Group[]> {
    return await this.get<Group[]>(`/users/${userId}/groups`);
  }
}
