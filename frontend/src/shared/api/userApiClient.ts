/**
 * User API Client
 */

import { V1ApiClient } from "./client";
import type { User } from "../types/domain";

/** Create user request */
export interface CreateUserRequest {
  name: string;
  email: string;
  role: string;
}

/** Update user request */
export interface UpdateUserRequest {
  name: string;
  email: string;
  role: string;
}

/** List users response */
interface ListUsersResponse {
  users: User[];
}

/** Delete user response */
interface DeleteUserResponse {
  message: string;
  userId: string;
}

export class UserApi extends V1ApiClient {
  constructor() {
    super("/users");
  }

  /**
   * Create a new user
   */
  async createUser(req: CreateUserRequest): Promise<User> {
    return this.post<User, CreateUserRequest>("", req);
  }

  /**
   * Get user by ID
   */
  async getUser(userId: string): Promise<User> {
    return this.get<User>(`/${userId}`);
  }

  /**
   * List all users
   */
  async listUsers(): Promise<User[]> {
    const response = await this.get<ListUsersResponse>("");
    return response.users;
  }

  /**
   * Update user
   */
  async updateUser(userId: string, req: UpdateUserRequest): Promise<User> {
    return this.put<User>(`/${userId}`, req);
  }

  /**
   * Delete user
   */
  async deleteUser(userId: string): Promise<DeleteUserResponse> {
    return this.delete<DeleteUserResponse>(`/${userId}`);
  }
}

// Export singleton instance for convenience
export const userApi = new UserApi();
