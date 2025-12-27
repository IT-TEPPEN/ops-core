import { NotificationInfo } from "./Notification";

export type NotificationsState = {
  notifications: NotificationInfo[];
};

export type NotificationsAction =
  | {
      type: "PUSH";
      payload: {
        title: string;
        message: string;
        type: "info" | "success" | "warning" | "error";
      };
    }
  | {
      type: "HIDE_TOP";
    }
  | {
      type: "POP";
    };
