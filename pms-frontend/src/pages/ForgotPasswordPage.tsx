// src/pages/ForgotPasswordPage.tsx

import { useForm } from "react-hook-form";
import type { SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";

// 1. Define the validation schema for a single email field
const forgotPasswordSchema = z.object({
  email: z.string()
    .min(1, "Email is required")
    .regex(/^[^\s@]+@[^\s@]+\.[^\s@]+$/, "Invalid email address"),
});

// 2. Infer the TypeScript type
type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export function ForgotPasswordPage() {
  // 3. Set up React Hook Form
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting, isSubmitSuccessful }, // Add isSubmitSuccessful
  } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
    mode: 'onSubmit',
    defaultValues: {
      email: ''
    }
  });

  // 4. Define the submit handler
  const onSubmit: SubmitHandler<ForgotPasswordFormData> = async (data) => {
    // API call to your backend to initiate the password reset process
    console.log("Submitting password reset request for:", data.email);
    await new Promise(resolve => setTimeout(resolve, 2000));
    console.log("Password reset email sent (simulation)!");
    // In a real app, you might show a success message based on API response
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-900">Forgot Your Password?</h2>
        
        {isSubmitSuccessful ? (
          <div className="text-center">
            <h3 className="text-lg font-medium text-gray-800">Check your email</h3>
            <p className="mt-2 text-sm text-gray-600">
              If an account with that email exists, we have sent instructions to reset your password.
            </p>
            <div className="mt-4">
              <a href="/login" className="font-medium text-blue-600 hover:text-blue-500">
                &larr; Back to Sign In
              </a>
            </div>
          </div>
        ) : (
          <>
            <p className="text-sm text-center text-gray-600">
              Enter the email address associated with your account, and we'll send you a link to reset your password.
            </p>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <Input
                id="email"
                type="email"
                label="Email Address"
                register={register("email")}
                error={errors.email}
                autoComplete="email"
              />
              <div>
                <Button type="submit" className="w-full" isLoading={isSubmitting}>
                  Send Reset Link
                </Button>
              </div>
            </form>
            <p className="text-sm text-center text-gray-600">
              Remember your password?{' '}
              <a href="/login" className="font-medium text-blue-600 hover:text-blue-500">
                Sign in
              </a>
            </p>
          </>
        )}
      </div>
    </div>
  );
}