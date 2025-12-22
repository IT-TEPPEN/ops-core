import React from "react";
import { TanstackClientProvider } from "./TanstackClientProvider";
import { RepositoryManagementAdapterProvider } from "./RepositoryManagementAdapterProvider";

export function DiProviders(props: { children: React.ReactNode }) {
  return (
    <TanstackClientProvider>
      <RepositoryManagementAdapterProvider>
        {props.children}
      </RepositoryManagementAdapterProvider>
    </TanstackClientProvider>
  );
}
