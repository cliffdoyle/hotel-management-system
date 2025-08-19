// src/router/PublicRoute.tsx
import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export function PublicRoute() {
  const { isAuthenticated } = useAuth();
  
  // If the user is authenticated, redirect them from public pages (like login) to the dashboard.
  return isAuthenticated ? <Navigate to="/dashboard" replace /> : <Outlet />;
}