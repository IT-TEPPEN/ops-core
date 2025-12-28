import { useEffect, useReducer } from "react";
import { useAuth } from "../app/hooks/useAuth";

const providerDisplayNames: Record<string, string> = {
  google: "Google",
  github: "GitHub",
  gitlab: "GitLab",
  microsoft: "Microsoft",
};

interface AccountSettingsState {
  loading: boolean;
  error: string | null;
}

type AccountSettingsAction =
  | { type: "SET_LOADING"; payload: boolean }
  | { type: "SET_ERROR"; payload: string | null };

function accountSettingsReducer(
  state: AccountSettingsState,
  action: AccountSettingsAction
): AccountSettingsState {
  switch (action.type) {
    case "SET_LOADING":
      return { ...state, loading: action.payload };
    case "SET_ERROR":
      return { ...state, error: action.payload };
    default:
      return state;
  }
}

export default function AccountSettingsPage() {
  const { user, identities, refreshIdentities, linkProvider, unlinkIdentity } =
    useAuth();
  const [state, dispatch] = useReducer(accountSettingsReducer, {
    loading: true,
    error: null,
  });

  useEffect(() => {
    const loadIdentities = async () => {
      try {
        await refreshIdentities();
      } catch {
        dispatch({
          type: "SET_ERROR",
          payload: "Failed to load linked accounts",
        });
      } finally {
        dispatch({ type: "SET_LOADING", payload: false });
      }
    };

    loadIdentities();
  }, [refreshIdentities]);

  const handleLinkProvider = async (provider: string) => {
    try {
      dispatch({ type: "SET_ERROR", payload: null });
      await linkProvider(provider);
    } catch {
      dispatch({
        type: "SET_ERROR",
        payload: `Failed to link ${
          providerDisplayNames[provider] || provider
        } account`,
      });
    }
  };

  const handleUnlinkIdentity = async (identityId: string) => {
    if (identities.length <= 1) {
      dispatch({
        type: "SET_ERROR",
        payload: "Cannot unlink your last identity",
      });
      return;
    }

    if (!confirm("Are you sure you want to unlink this account?")) {
      return;
    }

    try {
      dispatch({ type: "SET_ERROR", payload: null });
      await unlinkIdentity(identityId);
      await refreshIdentities();
    } catch {
      dispatch({ type: "SET_ERROR", payload: "Failed to unlink account" });
    }
  };

  const availableProviders = ["google", "github", "gitlab", "microsoft"].filter(
    (provider) => !identities.some((id) => id.provider === provider)
  );

  if (state.loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

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
                {user?.picture && (
                  <img
                    src={user.picture}
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

            {/* Error Message */}
            {state.error && (
              <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                {state.error}
              </div>
            )}

            {/* Linked Accounts */}
            <div className="mb-8">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Linked Accounts
              </h3>
              <div className="space-y-4">
                {identities.map((identity) => (
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
                    {identities.length > 1 && (
                      <button
                        onClick={() => handleUnlinkIdentity(identity.id)}
                        className="px-4 py-2 text-sm font-medium text-red-600 hover:text-red-700 border border-red-300 rounded-md hover:bg-red-50 transition-colors"
                      >
                        Unlink
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* Available Providers */}
            {availableProviders.length > 0 && (
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">
                  Link Additional Accounts
                </h3>
                <div className="space-y-3">
                  {availableProviders.map((provider) => (
                    <button
                      key={provider}
                      onClick={() => handleLinkProvider(provider)}
                      className="w-full flex items-center justify-between p-4 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                      <span className="text-sm font-medium text-gray-900">
                        {providerDisplayNames[provider] || provider}
                      </span>
                      <span className="text-sm text-blue-600">Link →</span>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
