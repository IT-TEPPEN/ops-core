import { ReactNode, useEffect, useMemo, useReducer, type Reducer } from "react";
import { NotificationsAction, NotificationsState } from "../types";
import { NotificationInfoImpl } from "../models";
import {
  NotificationsStateContext,
  NotificationsActionsContext,
} from "./NotificationContexts";

export type {
  NotificationsStateContextType,
  NotificationsActionContextType,
} from "./NotificationContexts";

// Reducer
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
    case "HIDE_TOP": {
      console.log("HIDE_TOP action dispatched");
      if (
        state.notifications.length === 0 ||
        state.notifications[0].getStatus() === "hidden"
      ) {
        console.log("No visible notifications to hide");
        return state;
      }

      console.log("Hiding top notification");

      const topNotification = state.notifications[0];
      const newTopNotification = topNotification.hide();
      return {
        notifications: [newTopNotification, ...state.notifications.slice(1)],
      };
    }
    case "POP": {
      const hidedAt = state.notifications[0]?.getHidedAt();

      if (
        state.notifications.length === 0 ||
        state.notifications[0].getStatus() === "visible" ||
        hidedAt === null ||
        new Date().getTime() - hidedAt.getTime() < 3000
      ) {
        return state;
      }

      return {
        notifications: state.notifications.slice(1),
      };
    }
  }
};

// Provider component
export function NotificationsProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(reducer, { notifications: [] });

  useEffect(() => {
    if (state.notifications.length === 0) {
      return;
    }

    const topNotification = state.notifications[0];
    if (topNotification.getStatus() === "hidden") {
      const timer = setTimeout(() => {
        for (const notification of state.notifications) {
          if (notification.getStatus() === "visible") {
            return;
          }
          dispatch({ type: "POP" });
        }
      }, 5000);

      return () => clearTimeout(timer);
    }
  }, [state.notifications]);

  const notificationActions = useMemo(
    () => ({
      push: (data: {
        title: string;
        message: string;
        type: "info" | "success" | "warning" | "error";
      }) => {
        dispatch({ type: "PUSH", payload: data });
      },
      dismiss: () => {
        console.log("dismiss called");
        dispatch({ type: "HIDE_TOP" });
      },
    }),
    []
  );

  return (
    <NotificationsStateContext.Provider value={state.notifications}>
      <NotificationsActionsContext.Provider value={notificationActions}>
        {children}
      </NotificationsActionsContext.Provider>
    </NotificationsStateContext.Provider>
  );
}
