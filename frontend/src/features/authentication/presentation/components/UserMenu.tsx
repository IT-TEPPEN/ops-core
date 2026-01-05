import { useEffect, useState } from "react";
import { useGetUserIdentityUsecase } from "../contexts";
import { useSignOut } from "../hooks/useSignOut";
import { useUserMenu } from "../hooks/useUserMenu";
import { UserAvatar } from "./UserAvatar";
import { UserMenuContent } from "./UserMenuContent";

export function UserMenu() {
  const { isOpen, toggleMenu, closeMenu, menuRef } = useUserMenu();
  const { signOut, isLoading } = useSignOut();
  const getUserIdentityUsecase = useGetUserIdentityUsecase();
  const [identity, setIdentity] = useState<{
    name: string;
    pictureUrl: string;
  } | null>(null);

  useEffect(() => {
    const fetchIdentity = async () => {
      try {
        const result = await getUserIdentityUsecase.execute();
        if (result) {
          setIdentity({
            name: result.name,
            pictureUrl: result.pictureUrl,
          });
        }
      } catch (error) {
        console.error("Failed to fetch user identity:", error);
      }
    };

    fetchIdentity();
  }, [getUserIdentityUsecase]);

  if (!identity) {
    // Don't render if no user is logged in
    return null;
  }

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={toggleMenu}
        disabled={isLoading}
        className="focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-full transition duration-150 ease-in-out disabled:opacity-50"
        aria-label="User menu"
      >
        <UserAvatar
          name={identity.name}
          pictureUrl={identity.pictureUrl}
          size="md"
        />
      </button>

      {isOpen && <UserMenuContent onSignOut={signOut} onClose={closeMenu} />}
    </div>
  );
}
