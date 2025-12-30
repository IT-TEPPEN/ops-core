import { useQuery } from "@tanstack/react-query";
import { useRepositoryQueryService } from "./useRepositoryQueryService";

/**
 * Hook for fetching repository list.
 * Uses React Query for caching and state management.
 * Following ADR 0019 - simple queries can use Query Service directly.
 */
export function useRepositoryList() {
  const queryService = useRepositoryQueryService();

  return useQuery({
    queryKey: ["repositories"],
    queryFn: () => queryService.list(),
  });
}
