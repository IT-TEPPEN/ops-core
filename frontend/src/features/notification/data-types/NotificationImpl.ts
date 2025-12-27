import { NotificationInfo } from "./Notification";

export class NotificationInfoImpl implements NotificationInfo {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly message: string,
    private readonly type: "info" | "success" | "warning" | "error",
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
  getCreatedAt(): Date {
    return this.createdAt;
  }
  getRead(): boolean {
    return this.read;
  }
}
