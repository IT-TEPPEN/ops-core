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
