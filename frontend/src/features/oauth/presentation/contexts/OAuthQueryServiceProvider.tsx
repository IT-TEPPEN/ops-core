import { ReactNode, useMemo } from "react";
import { HttpOAuthQueryService } from "../../infrastructure";
import { OAuthQueryServiceContext } from "./OAuthQueryServiceContext";

export function OAuthQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpOAuthQueryService();
  }, []);

  return (
    <OAuthQueryServiceContext.Provider value={service}>
      {children}
    </OAuthQueryServiceContext.Provider>
  );
}
