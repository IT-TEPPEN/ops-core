import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { RepositoryManagementAdapterProvider } from "./RepositoryManagementAdapterProvider";
import { NotificationsProvider } from "@/features/notification";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <RepositoryManagementAdapterProvider>
        <NotificationsProvider>{props.children}</NotificationsProvider>
      </RepositoryManagementAdapterProvider>
    </TanstackClientProvider>
  );
}
