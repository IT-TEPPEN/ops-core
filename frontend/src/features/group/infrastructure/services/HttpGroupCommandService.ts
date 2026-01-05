import type {
  GroupCommandService,
  CreateGroupRequest,
  UpdateGroupRequest,
  MemberRequest,
} from "../../application";
import type { Group } from "@/shared/types/domain";
import { V1ApiClient } from "@/shared/api/client";

/**
 * HTTP implementation of GroupCommandService.
 * Handles all write operations (POST/PUT/DELETE) for group management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpGroupCommandService
  extends V1ApiClient
  implements GroupCommandService
{
  constructor() {
    super("/groups");
  }

  async create(request: CreateGroupRequest): Promise<Group> {
    return await this.post<Group, CreateGroupRequest>("", request);
  }

  async update(groupId: string, request: UpdateGroupRequest): Promise<Group> {
    return await this.put<Group, UpdateGroupRequest>(`/${groupId}`, request);
  }

  async deleteGroup(groupId: string): Promise<void> {
    await this.delete<void>(`/${groupId}`);
  }

  async addMember(groupId: string, request: MemberRequest): Promise<Group> {
    return await this.post<Group, MemberRequest>(
      `/${groupId}/members`,
      request
    );
  }

  async removeMember(groupId: string, request: MemberRequest): Promise<Group> {
    return await this.deleteWithBody<Group, MemberRequest>(
      `/${groupId}/members`,
      request
    );
  }
}
