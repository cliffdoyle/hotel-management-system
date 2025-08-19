// src/pages/LoginPage.tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "../components/ui/Button";
import { useAuth } from "../hooks/useAuth";
import { loginUser } from "../api/authService";

export interface LoginFormData {
  email: string;
  password: string;
}

export function LoginPage() {
  const [apiError, setApiError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const auth = useAuth();
  const navigate = useNavigate();

  // Simple form state
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: ''
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setApiError(null);
    setIsSubmitting(true);

    console.log('Login form data being sent:', formData);

    // Basic validation
    if (!formData.email || !formData.email.includes('@')) {
      setApiError("Please enter a valid email address");
      setIsSubmitting(false);
      return;
    }
    if (!formData.password) {
      setApiError("Password is required");
      setIsSubmitting(false);
      return;
    }

    try {
      const response = await loginUser(formData);
      auth.login(response.token, response.user);
      navigate("/dashboard");
    } catch (error) {
      console.error('Login error:', error);
      if (error instanceof Error) {
        setApiError(error.message);
      } else {
        setApiError("An unknown error occurred.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-900">Sign in to your account</h2>
        {apiError && <div className="p-3 text-sm text-red-700 bg-red-100 border border-red-400 rounded-md">{apiError}</div>}
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email Address</label>
            <input
              id="email"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleInputChange}
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              autoComplete="email"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">Password</label>
            <input
              id="password"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleInputChange}
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              autoComplete="current-password"
            />
          </div>
          <div className="text-sm text-right">
            <Link to="/forgot-password" className="font-medium text-blue-600 hover:text-blue-500">Forgot your password?</Link>
          </div>
          <div>
            <Button type="submit" className="w-full" isLoading={isSubmitting}>Sign In</Button>
          </div>
        </form>
        <p className="text-sm text-center text-gray-600">Not a member?{' '}
          <Link to="/register" className="font-medium text-blue-600 hover:text-blue-500">Create an account</Link>
        </p>
      </div>
    </div>
  );
}