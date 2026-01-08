import { useMemo } from "react";
import { HttpDocumentQueryService } from "../../infrastructure";
import { DocumentQueryServiceContext } from "./DocumentQueryServiceContext";

export function DocumentQueryServiceProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpDocumentQueryService();
  }, []);

  return (
    <DocumentQueryServiceContext.Provider value={service}>
      {children}
    </DocumentQueryServiceContext.Provider>
  );
}
