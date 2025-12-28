export interface NotificationInfo {
  getId(): string;
  getTitle(): string;
  getMessage(): string;
  getType(): "info" | "success" | "warning" | "error";
  getStatus(): "visible" | "hidden";
  getCreatedAt(): Date;
  getRead(): boolean;
  hide(): NotificationInfo;
  getHidedAt(): Date | null;
}

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
