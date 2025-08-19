// src/pages/DashboardPage.tsx
import { Button } from "../components/ui/Button";
import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";

export function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  const displayName = user?.first_name && user?.last_name
    ? `${user.first_name} ${user.last_name}`
    : user?.first_name || user?.email || 'User';

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
        <div className="p-8 bg-white rounded-lg shadow-lg text-center">
          <h1 className="text-3xl font-bold">Welcome, {displayName}!</h1>
          <p className="mt-2 text-gray-600">Your email is: {user?.email}</p>
          <p className="mt-1 text-sm text-gray-500">Your roles: {user?.roles?.join(', ') || 'No roles assigned'}</p>
          <div className="mt-6">
            <Button onClick={handleLogout}>Logout</Button>
          </div>
        </div>
    </div>
  );
}