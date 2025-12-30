import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { AuthApiProvider } from "./AuthApiProvider";
import { NotificationsProvider } from "@/features/notification";
import {
  RepositoryQueryServiceProvider,
  RepositoryCommandServiceProvider,
} from "@/features/repository";
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
              <RepositoryCommandServiceProvider>
                <NotificationsProvider>{props.children}</NotificationsProvider>
              </RepositoryCommandServiceProvider>
            </RepositoryQueryServiceProvider>
          </OAuthCommandServiceProvider>
        </OAuthQueryServiceProvider>
      </AuthApiProvider>
    </TanstackClientProvider>
  );
}
