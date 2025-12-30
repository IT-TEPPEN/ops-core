import type { UserQueryService } from "../../application/services";
import type { User, Group } from "@/shared/types/domain";
import { V1ApiClient } from "@/shared/api/client";

/**
 * HTTP implementation of UserQueryService.
 * Handles all read operations (GET) for user management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpUserQueryService
  extends V1ApiClient
  implements UserQueryService
{
  constructor() {
    super("/users");
  }

  async listUsers(): Promise<User[]> {
    return await this.get<User[]>("");
  }

  async getUserGroups(userId: string): Promise<Group[]> {
    return await this.get<Group[]>(`/${userId}/groups`);
  }
}
