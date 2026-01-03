import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { NotificationsProvider } from "@/features/notification";
import { RepositoryQueryServiceProvider } from "@/features/repository";
import {
  OAuthCommandServiceProvider,
  OAuthQueryServiceProvider,
} from "@/features/oauth";
import { AuthenticationServiceProvider } from "@/features/authentication/infrastructure/contexts";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <AuthenticationServiceProvider>
        <OAuthQueryServiceProvider>
          <OAuthCommandServiceProvider>
            <RepositoryQueryServiceProvider>
              <NotificationsProvider>{props.children}</NotificationsProvider>
            </RepositoryQueryServiceProvider>
          </OAuthCommandServiceProvider>
        </OAuthQueryServiceProvider>
      </AuthenticationServiceProvider>
    </TanstackClientProvider>
  );
}
