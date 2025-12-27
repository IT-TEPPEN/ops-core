import { createContext } from "react";
import { NotificationsStateContextType } from "../data-types/NotificationsContext";

export const NotificationsStateContext =
  createContext<NotificationsStateContextType>([
    {
      getId: () => "sample-id-1",
      getTitle: () => "Sample Notification1",
      getMessage: () => "[1] This is a sample notification message.",
      getType: () => "info",
      getRead: () => false,
      getCreatedAt: () => new Date(),
    },
    {
      getId: () => "sample-id-2",
      getTitle: () => "Sample Notification2",
      getMessage: () => "[2] This is a sample notification message.",
      getType: () => "warning",
      getRead: () => false,
      getCreatedAt: () => new Date(),
    },
  ]);
