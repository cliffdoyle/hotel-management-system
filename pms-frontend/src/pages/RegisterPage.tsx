// src/pages/RegisterPage.tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "../components/ui/Button";
import { useAuth } from "../hooks/useAuth";
import { registerUser } from "../api/authService";

interface RegisterFormData {
    first_name: string;
    last_name: string;
    email: string;
    password: string;
}


export default function RegisterPage() {
    const [apiError, setApiError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const auth = useAuth();
    const navigate = useNavigate();

    // Simple form state
    const [formData, setFormData] = useState<RegisterFormData>({
        first_name: '',
        last_name: '',
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

        console.log('Form data being sent:', formData);

        // Basic validation
        if (!formData.first_name || formData.first_name.trim().length < 2) {
            setApiError("First name must be at least 2 characters");
            setIsSubmitting(false);
            return;
        }
        if (!formData.last_name || formData.last_name.trim().length < 2) {
            setApiError("Last name must be at least 2 characters");
            setIsSubmitting(false);
            return;
        }
        if (!formData.email || !formData.email.includes('@')) {
            setApiError("Please enter a valid email address");
            setIsSubmitting(false);
            return;
        }
        if (!formData.password || formData.password.length < 8) {
            setApiError("Password must be at least 8 characters");
            setIsSubmitting(false);
            return;
        }

        try {
            const response = await registerUser(formData);
            auth.login(response.token, response.user);
            navigate("/dashboard");
        } catch (error) {
            console.error('Registration error:', error);
            if (error instanceof Error) {
                setApiError(error.message);
            } else {
                setApiError(`Registration failed: ${JSON.stringify(error)}`);
            }
        } finally {
            setIsSubmitting(false);
        }
    };
  
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
            <h2 className="text-2xl font-bold text-center text-gray-900">Create a new account</h2>
            {apiError && <div className="p-3 text-sm text-red-700 bg-red-100 border border-red-400 rounded-md">{apiError}</div>}
            <form onSubmit={handleSubmit} className="space-y-6">
                <div>
                    <label htmlFor="first_name" className="block text-sm font-medium text-gray-700">First Name</label>
                    <input
                        id="first_name"
                        name="first_name"
                        type="text"
                        value={formData.first_name}
                        onChange={handleInputChange}
                        className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        autoComplete="given-name"
                    />
                </div>
                <div>
                    <label htmlFor="last_name" className="block text-sm font-medium text-gray-700">Last Name</label>
                    <input
                        id="last_name"
                        name="last_name"
                        type="text"
                        value={formData.last_name}
                        onChange={handleInputChange}
                        className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        autoComplete="family-name"
                    />
                </div>
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
                        autoComplete="new-password"
                    />
                </div>
                <div>
                    <Button type="submit" className="w-full" isLoading={isSubmitting}>
                        Create Account
                    </Button>
                </div>
            </form>
            <p className="text-sm text-center text-gray-600">Already have an account?{' '}
                <Link to="/login" className="font-medium text-blue-600 hover:text-blue-500">Sign in</Link>
            </p>
        </div>
      </div>
    );
}