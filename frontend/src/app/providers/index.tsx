import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { RepositoryManagementAdapterProvider } from "./RepositoryManagementAdapterProvider";
import { AuthApiProvider } from "./AuthApiProvider";
import { NotificationsProvider } from "@/features/notification";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <AuthApiProvider>
        <RepositoryManagementAdapterProvider>
          <NotificationsProvider>{props.children}</NotificationsProvider>
        </RepositoryManagementAdapterProvider>
      </AuthApiProvider>
    </TanstackClientProvider>
  );
}
