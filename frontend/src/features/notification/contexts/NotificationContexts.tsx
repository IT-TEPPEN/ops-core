import { createContext } from "react";
import { NotificationInfo } from "../types";

// Context types
export type NotificationsStateContextType = NotificationInfo[];

export type NotificationsActionContextType = {
  push: (data: {
    title: string;
    message: string;
    type: "info" | "success" | "warning" | "error";
  }) => void;
  dismiss: () => void;
};

// Create contexts
export const NotificationsStateContext =
  createContext<NotificationsStateContextType>([]);

export const NotificationsActionsContext =
  createContext<NotificationsActionContextType>({
    push: () => {
      throw new Error("Function not implemented.");
    },
    dismiss: () => {
      throw new Error("Function not implemented.");
    },
  });
