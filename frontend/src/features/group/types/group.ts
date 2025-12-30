/**
 * Group feature type definitions.
 * Following Pattern 1 - Simplified structure.
 */

/**
 * Group update request.
 */
export interface GroupUpdateRequest {
  name: string;
  description: string;
}

/**
 * Member management request.
 */
export interface MemberRequest {
  user_id: string;
}

/**
 * Group detail state for UI.
 */
export interface GroupDetailState {
  group: any | null; // Using shared types from @/shared/types/domain
  members: any[];
  isLoading: boolean;
  isEditing: boolean;
  showMemberSelector: boolean;
  error: string | null;
}

/**
 * Group detail actions.
 */
export type GroupDetailAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; group: any; members: any[] }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "SET_EDITING"; isEditing: boolean }
  | { type: "SET_MEMBER_SELECTOR"; show: boolean }
  | { type: "UPDATE_SUCCESS"; group: any }
  | { type: "UPDATE_MEMBERS"; members: any[] }
  | { type: "SET_ERROR"; error: string };
