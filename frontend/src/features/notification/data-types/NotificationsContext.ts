import { NotificationInfo } from "./Notification";

export type NotificationsStateContextType = NotificationInfo[];

type PushFunction = (data: {
  title: string;
  message: string;
  type: "info" | "success" | "warning" | "error";
}) => void;

export type NotificationsActionContextType = {
  push: PushFunction;
};
