/**
 * User API utilities
 */

import { get, post, put, del } from "./client";
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
export interface ListUsersResponse {
  users: User[];
}

/**
 * Create a new user
 */
export async function createUser(
  req: CreateUserRequest,
  signal?: AbortSignal
): Promise<User> {
  return post<User>("/users", req, signal);
}

/**
 * Get user by ID
 */
export async function getUser(
  userId: string,
  signal?: AbortSignal
): Promise<User> {
  return get<User>(`/users/${userId}`, signal);
}

/**
 * List all users
 */
export async function listUsers(signal?: AbortSignal): Promise<User[]> {
  const response = await get<ListUsersResponse>("/users", signal);
  return response.users;
}

/**
 * Update user
 */
export async function updateUser(
  userId: string,
  req: UpdateUserRequest,
  signal?: AbortSignal
): Promise<User> {
  return put<User>(`/users/${userId}`, req, signal);
}

/**
 * Delete user
 */
export async function deleteUser(
  userId: string,
  signal?: AbortSignal
): Promise<{ message: string; userId: string }> {
  return del<{ message: string; userId: string }>(`/users/${userId}`, signal);
}
