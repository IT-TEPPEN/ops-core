import { TypedParamsGuard } from "@/shared/components";
import { Page } from "@/shared/types/Page";

export function DocumentGuard(props: { element: Page<"docId"> }) {
  return (
    <TypedParamsGuard
      pathParams={[{ name: "docId", required: true }]}
      queryParams={undefined}
      fallbackPath="/documents"
      errorMessage="Invalid document ID."
      element={props.element}
    />
  );
}
