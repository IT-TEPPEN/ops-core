import React, { JSX } from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { NotificationsProvider } from "@/features/notification";
import { RepositoryQueryServiceProvider } from "@/features/repository";
import {
  OAuthCommandServiceProvider,
  OAuthQueryServiceProvider,
} from "@/features/oauth";
import { AuthenticationDiProvider } from "@/features/authentication/presentation/contexts";
import {
  DocumentCommandServiceProvider,
  DocumentQueryServiceProvider,
} from "@/features/document";

const Providers: ((props: { children: React.ReactNode }) => JSX.Element)[] = [
  TanstackClientProvider,
  AuthenticationDiProvider,
  OAuthQueryServiceProvider,
  OAuthCommandServiceProvider,
  DocumentQueryServiceProvider,
  DocumentCommandServiceProvider,
  RepositoryQueryServiceProvider,
  NotificationsProvider,
];

export function DiProviders(props: { children: React.ReactNode }) {
  return Providers.reduceRight(
    (children, Provider) => <Provider>{children}</Provider>,
    props.children
  );
}
