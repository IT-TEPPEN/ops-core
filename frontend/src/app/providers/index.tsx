import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { AuthApiProvider } from "./AuthApiProvider";
import { NotificationsProvider } from "@/features/notification";
import {
  RepositoryQueryServiceProvider,
  RepositoryCommandServiceProvider,
} from "@/features/repository";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <AuthApiProvider>
        <RepositoryQueryServiceProvider>
          <RepositoryCommandServiceProvider>
            <NotificationsProvider>{props.children}</NotificationsProvider>
          </RepositoryCommandServiceProvider>
        </RepositoryQueryServiceProvider>
      </AuthApiProvider>
    </TanstackClientProvider>
  );
}
