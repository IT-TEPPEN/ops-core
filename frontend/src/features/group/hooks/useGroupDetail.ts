import { useReducer, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { Group, User } from "@/shared/types/domain";
import {
  getGroup,
  listUsers,
  updateGroup,
  deleteGroup,
  addMember,
  removeMember,
} from "@/shared/api";

interface GroupDetailState {
  group: Group | null;
  members: User[];
  isLoading: boolean;
  isEditing: boolean;
  showMemberSelector: boolean;
  error: string | null;
}

type GroupDetailAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; group: Group; members: User[] }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "SET_EDITING"; isEditing: boolean }
  | { type: "SET_MEMBER_SELECTOR"; show: boolean }
  | { type: "UPDATE_SUCCESS"; group: Group }
  | { type: "UPDATE_MEMBERS"; members: User[] }
  | { type: "SET_ERROR"; error: string };

function groupDetailReducer(
  state: GroupDetailState,
  action: GroupDetailAction
): GroupDetailState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, error: null };

    case "FETCH_SUCCESS":
      return {
        ...state,
        group: action.group,
        members: action.members,
        isLoading: false,
      };

    case "FETCH_ERROR":
      return {
        ...state,
        isLoading: false,
        error: action.error,
        group: null,
      };

    case "SET_EDITING":
      return { ...state, isEditing: action.isEditing };

    case "SET_MEMBER_SELECTOR":
      return { ...state, showMemberSelector: action.show };

    case "UPDATE_SUCCESS":
      return {
        ...state,
        group: action.group,
        isEditing: false,
      };

    case "UPDATE_MEMBERS":
      return {
        ...state,
        members: action.members,
      };

    case "SET_ERROR":
      return { ...state, error: action.error };

    default:
      return state;
  }
}

const initialState: GroupDetailState = {
  group: null,
  members: [],
  isLoading: true,
  isEditing: false,
  showMemberSelector: false,
  error: null,
};

export function useGroupDetail(groupId: string | undefined) {
  const [state, dispatch] = useReducer(groupDetailReducer, initialState);
  const navigate = useNavigate();

  const fetchGroup = useCallback(async () => {
    if (!groupId) return;

    dispatch({ type: "FETCH_START" });
    try {
      const data = await getGroup(groupId);
      const allUsers = await listUsers();
      const groupMembers = allUsers.filter((user) =>
        data.member_ids.includes(user.id)
      );
      dispatch({ type: "FETCH_SUCCESS", group: data, members: groupMembers });
    } catch (err) {
      dispatch({
        type: "FETCH_ERROR",
        error: err instanceof Error ? err.message : "Failed to load group",
      });
    }
  }, [groupId]);

  useEffect(() => {
    fetchGroup();
  }, [fetchGroup]);

  const handleUpdate = useCallback(
    async (data: { name: string; description: string }) => {
      if (!groupId) return;

      try {
        const updated = await updateGroup(groupId, data);
        dispatch({ type: "UPDATE_SUCCESS", group: updated });
      } catch (err) {
        throw err; // Let the form handle the error
      }
    },
    [groupId]
  );

  const handleDelete = useCallback(async () => {
    if (
      !groupId ||
      !window.confirm("Are you sure you want to delete this group?")
    ) {
      return;
    }

    try {
      await deleteGroup(groupId);
      navigate("/groups");
    } catch (err) {
      dispatch({
        type: "SET_ERROR",
        error: err instanceof Error ? err.message : "Failed to delete group",
      });
    }
  }, [groupId, navigate]);

  const handleAddMember = useCallback(
    async (userId: string) => {
      if (!groupId) return;

      try {
        const updated = await addMember(groupId, { user_id: userId });
        dispatch({ type: "UPDATE_SUCCESS", group: updated });
        dispatch({ type: "SET_MEMBER_SELECTOR", show: false });

        // Update member list
        const allUsers = await listUsers();
        const groupMembers = allUsers.filter((user) =>
          updated.member_ids.includes(user.id)
        );
        dispatch({ type: "UPDATE_MEMBERS", members: groupMembers });
      } catch (err) {
        throw err; // Let the selector handle the error
      }
    },
    [groupId]
  );

  const handleRemoveMember = useCallback(
    async (userId: string) => {
      if (
        !groupId ||
        !window.confirm("Are you sure you want to remove this member?")
      ) {
        return;
      }

      try {
        const updated = await removeMember(groupId, { user_id: userId });
        dispatch({ type: "UPDATE_SUCCESS", group: updated });

        // Update member list by filtering out removed user
        const updatedMembers = state.members.filter(
          (member) => member.id !== userId
        );
        dispatch({ type: "UPDATE_MEMBERS", members: updatedMembers });
      } catch (err) {
        dispatch({
          type: "SET_ERROR",
          error: err instanceof Error ? err.message : "Failed to remove member",
        });
      }
    },
    [groupId, state.members]
  );

  const setEditing = useCallback((isEditing: boolean) => {
    dispatch({ type: "SET_EDITING", isEditing });
  }, []);

  const setShowMemberSelector = useCallback((show: boolean) => {
    dispatch({ type: "SET_MEMBER_SELECTOR", show });
  }, []);

  return {
    group: state.group,
    members: state.members,
    isLoading: state.isLoading,
    isEditing: state.isEditing,
    showMemberSelector: state.showMemberSelector,
    error: state.error,
    handleUpdate,
    handleDelete,
    handleAddMember,
    handleRemoveMember,
    setEditing,
    setShowMemberSelector,
  };
}
