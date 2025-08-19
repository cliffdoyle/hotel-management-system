// src/App.tsx
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { LoginPage } from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import { ForgotPasswordPage } from "./pages/ForgotPasswordPage";
import { ProtectedRoute } from "./router/ProtectedRoute";
import { DashboardPage } from "./pages/DashboardPage";
import { PublicRoute } from "./router/PublicRoute";
import { AuthProvider } from "./context/AuthContext";

function App() {
  return (
    <Routes>
      {/* Public routes that redirect if the user is authenticated */}
      <Route path="/" element={<PublicRoute />}>
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="forgot-password" element={<ForgotPasswordPage />} />
        <Route index element={<Navigate to="/login" replace />} />
      </Route>

      {/* Protected routes that require authentication */}
      <Route path="/" element={<ProtectedRoute />}>
        <Route path="dashboard" element={<DashboardPage />} />
      </Route>

      {/* Fallback route: redirects to login for unauthenticated users */}
      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  );
}

function AppWrapper() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <App />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default AppWrapper;