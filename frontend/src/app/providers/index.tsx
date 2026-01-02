import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { AuthApiProvider } from "./AuthApiProvider";
import { NotificationsProvider } from "@/features/notification";
import { RepositoryQueryServiceProvider } from "@/features/repository";
import {
  OAuthCommandServiceProvider,
  OAuthQueryServiceProvider,
} from "@/features/oauth";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <AuthApiProvider>
        <OAuthQueryServiceProvider>
          <OAuthCommandServiceProvider>
            <RepositoryQueryServiceProvider>
              <NotificationsProvider>{props.children}</NotificationsProvider>
            </RepositoryQueryServiceProvider>
          </OAuthCommandServiceProvider>
        </OAuthQueryServiceProvider>
      </AuthApiProvider>
    </TanstackClientProvider>
  );
}
