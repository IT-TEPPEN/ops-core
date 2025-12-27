import {
  NotificationsActionsContext,
  NotificationsStateContext,
} from "../contexts";
import { useNotificationHandler } from "../hooks/NotificationHandler";

export function NotificationsProvider(props: { children: React.ReactNode }) {
  const { state, actions } = useNotificationHandler();
  return (
    <NotificationsStateContext.Provider value={state}>
      <NotificationsActionsContext.Provider value={actions}>
        {props.children}
      </NotificationsActionsContext.Provider>
    </NotificationsStateContext.Provider>
  );
}
