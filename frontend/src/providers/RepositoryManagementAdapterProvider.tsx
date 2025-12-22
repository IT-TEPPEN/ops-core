import { RepositoryApi } from "../api/repositoryApi";
import { RepositoryManagementAdapterContext } from "../features/repository/contexts";

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
