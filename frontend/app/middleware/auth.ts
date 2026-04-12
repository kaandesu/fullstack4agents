import { useAuthManager } from "~/stores/auth-manager";
import { storeToRefs } from "pinia";
export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) return;

  const auth = useAuthManager();
  const { isAuthenticated, initialized } = storeToRefs(auth);

  // Routes that don't require authentication
  const publicRoutes = ["/", "/auth/sign-in", "/auth/sign-up"];

  // On first navigation to a protected route, try to refresh the session.
  // This handles page reloads where the PocketBase token is still valid.
  if (!initialized.value && !publicRoutes.includes(to.path)) {
    initialized.value = await auth.refreshToken();
  }

  // Not authenticated → redirect to sign-in
  if (!isAuthenticated.value && !publicRoutes.includes(to.path)) {
    return navigateTo("/auth/sign-in");
  }

  // Already authenticated → redirect away from auth pages
  if (isAuthenticated.value && to.path.startsWith("/auth")) {
    return navigateTo("/dashboard"); // change this to your post-login route
  }

  // Add any post-auth setup here (e.g. load user workspace, onboarding check)
});
