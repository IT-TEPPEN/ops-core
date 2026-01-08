import { useMemo } from "react";
import { HttpDocumentCommandService } from "../../infrastructure";
import { DocumentCommandServiceContext } from "./DocumentCommandServiceContext";

export function DocumentCommandServiceProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpDocumentCommandService();
  }, []);

  return (
    <DocumentCommandServiceContext.Provider value={service}>
      {children}
    </DocumentCommandServiceContext.Provider>
  );
}
