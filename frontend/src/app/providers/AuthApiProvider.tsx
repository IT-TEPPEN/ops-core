import { AuthApiContext } from "@/app/contexts/AuthApiContext";
import { AuthApi } from "@/shared/api/authApi";

export function AuthApiProvider(props: { children: React.ReactNode }) {
  const authApi = new AuthApi();

  return (
    <AuthApiContext.Provider value={authApi}>
      {props.children}
    </AuthApiContext.Provider>
  );
}
