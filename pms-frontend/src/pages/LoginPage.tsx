// src/pages/LoginPage.tsx

import  { useForm} from "react-hook-form";
import type { SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";

// 1. Define the validation schema with Zod
const loginSchema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters long" }),
});

// 2. Infer the TypeScript type from the Zod schema
type LoginFormData = z.infer<typeof loginSchema>;

export function LoginPage() {
  // 3. Set up React Hook Form
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  // 4. Define the submit handler
  const onSubmit: SubmitHandler<LoginFormData> = async (data) => {
    // This is where you'll call your backend API
    console.log("Form data:", data);
    // Simulate network delay for loading state
    await new Promise(resolve => setTimeout(resolve, 2000));
    console.log("Pretend login successful!");
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-900">Sign in to your account</h2>
        
        {/* 5. Use the handleSubmit wrapper */}
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            id="email"
            type="email"
            label="Email Address"
            // 6. Register the input and connect it to the form state
            register={register("email")}
            error={errors.email}
            autoComplete="email"
          />

          <Input
            id="password"
            type="password"
            label="Password"
            register={register("password")}
            error={errors.password}
            autoComplete="current-password"
          />

          <div className="text-sm text-right">
            <a href="/forgot-password" className="font-medium text-blue-600 hover:text-blue-500">
              Forgot your password?
            </a>
          </div>

          <div>
            <Button type="submit" className="w-full" isLoading={isSubmitting}>
              Sign In
            </Button>
          </div>
        </form>

        <p className="text-sm text-center text-gray-600">
          Not a member?{' '}
          <a href="/register" className="font-medium text-blue-600 hover:text-blue-500">
            Create an account
          </a>
        </p>
      </div>
    </div>
  );
}