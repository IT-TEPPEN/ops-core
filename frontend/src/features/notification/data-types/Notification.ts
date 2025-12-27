export interface NotificationInfo {
  getId(): string;
  getTitle(): string;
  getMessage(): string;
  getType(): "info" | "success" | "warning" | "error";
  getCreatedAt(): Date;
  getRead(): boolean;
}
