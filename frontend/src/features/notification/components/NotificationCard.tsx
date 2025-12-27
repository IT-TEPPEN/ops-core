import { useNotifications } from "../hooks";

export function NotificationCard() {
  const { state: notifications, actions } = useNotifications();

  return (
    <>
      {notifications.map((notification) => {
        let bgColor = "bg-white";

        switch (notification.getType()) {
          case "info":
            bgColor = "bg-blue-100";
            break;
          case "success":
            bgColor = "bg-green-100";
            break;
          case "warning":
            bgColor = "bg-yellow-100";
            break;
          case "error":
            bgColor = "bg-red-100";
            break;
        }

        return (
          <div
            key={notification.getId()}
            className={`relative py-4 px-8 shadow rounded-2xl duration-300 ${bgColor} ${
              notification.getStatus() === "hidden"
                ? "opacity-0 pointer-events-none"
                : "opacity-100"
            }`}
          >
            <div key={notification.getId()} className="mb-4 last:mb-0">
              <h3 className="text-lg font-semibold mb-1">
                {notification.getTitle()}
              </h3>
              <p className="text-gray-700 dark:text-gray-300">
                {notification.getMessage()}
              </p>
            </div>

            <div className="absolute top-2 right-2">
              <button
                onClick={() => {
                  console.log("Dismissing notification", notification.getId());
                  actions.dismiss();
                }}
                className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 cursor-pointer"
              >
                &times;
              </button>
            </div>
          </div>
        );
      })}
    </>
  );
}
