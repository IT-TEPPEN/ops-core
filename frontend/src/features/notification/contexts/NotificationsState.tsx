import { createContext } from "react";
import { NotificationsStateContextType } from "../data-types/NotificationsContext";

export const NotificationsStateContext =
  createContext<NotificationsStateContextType>([]);
