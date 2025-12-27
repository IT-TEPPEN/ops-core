import { createContext } from "react";
import { NotificationsActionContextType } from "../data-types/NotificationsContext";

export const NotificationsActionsContext =
  createContext<NotificationsActionContextType>({
    push: () => {
      throw new Error("Function not implemented.");
    },
  });
