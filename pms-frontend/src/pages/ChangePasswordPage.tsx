// src/pages/ChangePasswordPage.tsx
import { useForm } from "react-hook-form";
import { useState } from "react";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";
import { useAuth } from "../hooks/useAuth";

interface ChangePasswordFormData {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

export function ChangePasswordPage() {
  const { user } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    reset,
  } = useForm<ChangePasswordFormData>({
    mode: 'onSubmit',
    defaultValues: {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  });

  const newPassword = watch("newPassword");

  const onSubmit = async (data: ChangePasswordFormData) => {
    if (data.newPassword !== data.confirmPassword) {
      setMessage({ type: 'error', text: 'New passwords do not match' });
      return;
    }

    setIsLoading(true);
    setMessage(null);

    try {
      // TODO: Implement API call to change password
      // const response = await authService.changePassword({
      //   currentPassword: data.currentPassword,
      //   newPassword: data.newPassword,
      // });

      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setMessage({ type: 'success', text: 'Password changed successfully!' });
      reset();
    } catch (error) {
      setMessage({ 
        type: 'error', 
        text: error instanceof Error ? error.message : 'Failed to change password' 
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-lg">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Change Password</h1>
          <p className="mt-2 text-sm text-gray-600">
            Update your password for {user?.email}
          </p>
        </div>

        {message && (
          <div className={`mb-4 p-3 rounded-md ${
            message.type === 'success' 
              ? 'bg-green-50 text-green-800 border border-green-200' 
              : 'bg-red-50 text-red-800 border border-red-200'
          }`}>
            {message.text}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            id="currentPassword"
            label="Current Password"
            type="password"
            register={register("currentPassword", {
              required: "Current password is required",
            })}
            error={errors.currentPassword}
            autoComplete="current-password"
          />

          <Input
            id="newPassword"
            label="New Password"
            type="password"
            register={register("newPassword", {
              required: "New password is required",
              minLength: {
                value: 8,
                message: "Password must be at least 8 characters long",
              },
              pattern: {
                value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
                message: "Password must contain at least one uppercase letter, one lowercase letter, and one number",
              },
            })}
            error={errors.newPassword}
            autoComplete="new-password"
          />

          <Input
            id="confirmPassword"
            label="Confirm New Password"
            type="password"
            register={register("confirmPassword", {
              required: "Please confirm your new password",
              validate: (value) =>
                value === newPassword || "Passwords do not match",
            })}
            error={errors.confirmPassword}
            autoComplete="new-password"
          />

          <Button
            type="submit"
            isLoading={isLoading}
            className="w-full"
          >
            Change Password
          </Button>
        </form>

        <div className="mt-6 text-center">
          <a
            href="/dashboard"
            className="text-sm text-blue-600 hover:text-blue-500 transition-colors"
          >
            Back to Dashboard
          </a>
        </div>
      </div>
    </div>
  );
}
