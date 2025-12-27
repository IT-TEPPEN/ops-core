import { useContext } from "react";
import {
  NotificationsActionsContext,
  NotificationsStateContext,
} from "../contexts";

export function useNotificationsState() {
  return useContext(NotificationsStateContext);
}

export function useNotificationsActions() {
  return useContext(NotificationsActionsContext);
}

export function useNotifications() {
  return {
    state: useNotificationsState(),
    actions: useNotificationsActions(),
  };
}
