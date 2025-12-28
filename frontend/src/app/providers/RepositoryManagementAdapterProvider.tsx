import { RepositoryManagementAdapterContext } from "@/features/repository";
import { RepositoryApi } from "@/shared/api/repositoryApi";

export function RepositoryManagementAdapterProvider(props: {
  children: React.ReactNode;
}) {
  const adapter = new RepositoryApi();

  return (
    <RepositoryManagementAdapterContext.Provider value={adapter}>
      {props.children}
    </RepositoryManagementAdapterContext.Provider>
  );
}
