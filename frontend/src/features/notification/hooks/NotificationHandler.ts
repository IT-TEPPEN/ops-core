import { useMemo, useReducer, type Reducer } from "react";
import {
  NotificationsAction,
  NotificationsState,
  NotificationInfoImpl,
} from "../data-types";
import { ReturnReducerHooks } from "../../../types/ReturnReducerHooks";
import {
  NotificationsActionContextType,
  NotificationsStateContextType,
} from "../data-types/NotificationsContext";

const reducer: Reducer<NotificationsState, NotificationsAction> = (
  state,
  action
) => {
  switch (action.type) {
    case "PUSH": {
      const newNotification = NotificationInfoImpl.new({
        title: action.payload.title,
        message: action.payload.message,
        type: action.payload.type,
      });

      return {
        notifications: [...state.notifications, newNotification],
      };
    }
  }
};

export function useNotificationHandler(): ReturnReducerHooks<
  NotificationsStateContextType,
  NotificationsActionContextType
> {
  const [state, dispatch] = useReducer(reducer, { notifications: [] });

  const notificationActions = useMemo(
    () => ({
      push: (data: {
        title: string;
        message: string;
        type: "info" | "success" | "warning" | "error";
      }) => {
        dispatch({ type: "PUSH", payload: data });
      },
    }),
    [dispatch]
  );

  return {
    state: state.notifications,
    actions: notificationActions,
  };
}
