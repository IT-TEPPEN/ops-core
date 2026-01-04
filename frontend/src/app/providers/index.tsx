import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { NotificationsProvider } from "@/features/notification";
import { RepositoryQueryServiceProvider } from "@/features/repository";
// import {
//   OAuthCommandServiceProvider,
//   OAuthQueryServiceProvider,
// } from "@/features/oauth";
import { AuthenticationDiProvider } from "@/features/authentication/presentation/contexts";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <AuthenticationDiProvider>
        {/* <OAuthQueryServiceProvider>
          <OAuthCommandServiceProvider> */}
        <RepositoryQueryServiceProvider>
          <NotificationsProvider>{props.children}</NotificationsProvider>
        </RepositoryQueryServiceProvider>
        {/* </OAuthCommandServiceProvider>
        </OAuthQueryServiceProvider> */}
      </AuthenticationDiProvider>
    </TanstackClientProvider>
  );
}
