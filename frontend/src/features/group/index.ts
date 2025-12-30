// Group feature public API
// Following ADR 0021 - Pattern 1 (Simplified structure)

// Types
export type {
  GroupUpdateRequest,
  MemberRequest,
  GroupDetailState,
  GroupDetailAction,
} from "./types";

// Services
export type {
  GroupQueryService,
  GroupCommandService,
  CreateGroupRequest,
  UpdateGroupRequest,
} from "./application/services";

// Contexts
export {
  GroupQueryServiceProvider,
  useGroupQueryService,
  GroupCommandServiceProvider,
  useGroupCommandService,
} from "./presentation/contexts";

// Hooks
export { useGroupDetail } from "./hooks";

// Components
export { GroupInfoCard } from "./components/GroupInfoCard";
export { MembersPanel } from "./components/MembersPanel";
