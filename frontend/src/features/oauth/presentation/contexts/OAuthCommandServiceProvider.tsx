import { ReactNode, useMemo } from "react";
import { HttpOAuthCommandService } from "../../infrastructure";
import { OAuthCommandServiceContext } from "./OAuthCommandServiceContext";

export function OAuthCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpOAuthCommandService();
  }, []);

  return (
    <OAuthCommandServiceContext.Provider value={service}>
      {children}
    </OAuthCommandServiceContext.Provider>
  );
}
