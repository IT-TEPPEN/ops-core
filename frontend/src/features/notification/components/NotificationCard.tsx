import { useNotificationsState } from "../hooks";

export function NotificationCard() {
  const notifications = useNotificationsState();

  return (
    <div className="py-4 px-8 bg-white shadow rounded-2xl">
      {notifications.map((notification) => (
        <div key={notification.getId()} className="mb-4 last:mb-0">
          <h3 className="text-lg font-semibold mb-1">
            {notification.getTitle()}
          </h3>
          <p className="text-gray-700 dark:text-gray-300">
            {notification.getMessage()}
          </p>
        </div>
      ))}
    </div>
  );
}
