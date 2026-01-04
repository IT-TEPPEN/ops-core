import { useAuthenticationService } from "@/features/authentication/infrastructure/contexts";
import { useQuery } from "@tanstack/react-query";

const providerDisplayNames: Record<string, string> = {
  google: "Google",
  github: "GitHub",
  gitlab: "GitLab",
  microsoft: "Microsoft",
};
export default function AccountSettingsPage() {
  const authenticationService = useAuthenticationService();
  const query0 = useQuery({
    queryKey: ["identities"],
    queryFn: async () => {
      return authenticationService.listIdentities();
    },
  });
  const query1 = useQuery({
    queryKey: ["userProfile"],
    queryFn: async () => {
      return authenticationService.getIdentity();
    },
  });

  if (query0.isLoading || query1.isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (query0.isError || query1.isError) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {query0.isError && (
            <div>Error loading identities: {`${query0.error}`}</div>
          )}
          {query1.isError && (
            <div>Error loading user profile: {`${query1.error}`}</div>
          )}
        </div>
      </div>
    );
  }

  const identities = query0.data;
  const user = query1.data;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">
              Account Settings
            </h2>

            {/* User Profile */}
            <div className="mb-8">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Profile
              </h3>
              <div className="flex items-center space-x-4">
                {user?.pictureUrl && (
                  <img
                    src={user.pictureUrl}
                    alt={user.name}
                    className="h-16 w-16 rounded-full"
                  />
                )}
                <div>
                  <p className="text-lg font-medium text-gray-900">
                    {user?.name}
                  </p>
                  <p className="text-sm text-gray-500">{user?.email}</p>
                </div>
              </div>
            </div>

            {/* Linked Accounts */}
            <div className="mb-8">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Linked Accounts
              </h3>
              <div className="space-y-4">
                {identities?.identities.map((identity) => (
                  <div
                    key={identity.id}
                    className="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
                  >
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {providerDisplayNames[identity.provider] ||
                          identity.provider}
                      </p>
                      <p className="text-sm text-gray-500">{identity.email}</p>
                      <p className="text-xs text-gray-400">
                        Linked:{" "}
                        {new Date(identity.linkedAt).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
