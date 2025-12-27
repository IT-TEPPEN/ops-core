import { NotificationInfo } from "./Notification";

export class NotificationInfoImpl implements NotificationInfo {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly message: string,
    private readonly type: "info" | "success" | "warning" | "error",
    private status: "visible" | "hidden",
    private hidedAt: Date | null = null,
    private readonly createdAt: Date,
    private readonly read: boolean
  ) {}

  static new(data: {
    title: string;
    message: string;
    type: "info" | "success" | "warning" | "error";
  }): NotificationInfoImpl {
    return new NotificationInfoImpl(
      crypto.randomUUID(),
      data.title,
      data.message,
      data.type,
      "visible",
      null,
      new Date(),
      false
    );
  }

  getId(): string {
    return this.id;
  }
  getTitle(): string {
    return this.title;
  }
  getMessage(): string {
    return this.message;
  }
  getType(): "info" | "success" | "warning" | "error" {
    return this.type;
  }
  getStatus(): "visible" | "hidden" {
    return this.status;
  }
  getCreatedAt(): Date {
    return this.createdAt;
  }
  getRead(): boolean {
    return this.read;
  }
  getHidedAt(): Date | null {
    return this.hidedAt;
  }

  hide(): NotificationInfo {
    return new NotificationInfoImpl(
      this.id,
      this.title,
      this.message,
      this.type,
      "hidden",
      new Date(),
      this.createdAt,
      this.read
    );
  }
}
