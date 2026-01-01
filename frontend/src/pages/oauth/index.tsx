import { Route } from "react-router-dom";
import OAuthCallbackPage from "../OAuthCallbackPage";
import { OAuthConnection } from "@/features/oauth";

export function OAuthRoutes() {
  return (
    <Route path="oauth">
      <Route path="authorize" element={<OAuthConnection />} />
      <Route path="callback" element={<OAuthCallbackPage />} />
    </Route>
  );
}
