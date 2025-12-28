import { TypedParamsGuard } from "@/shared/components";
import { Page } from "@/shared/types/Page";

export function RepositoryIdGuard(props: { element: Page<"repoId"> }) {
  return (
    <TypedParamsGuard
      pathParams={[{ name: "repoId", required: true }]}
      queryParams={undefined}
      fallbackPath="/repositories"
      errorMessage="Invalid repository ID."
      element={props.element}
    />
  );
}
