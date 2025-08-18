import { useForm } from "react-hook-form";
import type { SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";

// 1. Define the validation schema with Zod, including password confirmation
const registerSchema = z.object({
  name: z.string().min(2, { message: "Name must be at least 2 characters long" }),
  email: z.string().email({ message: "Invalid email address" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
  confirmPassword: z.string().min(6, { message: "Password must be at least 6 characters" }),
})
.refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"], // Set the error on the confirmPassword field
});

// 2. Infer the TypeScript type from the Zod schema
type RegisterFormData = z.infer<typeof registerSchema>;

export function RegisterPage() {
  // 3. Set up React Hook Form
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  // 4. Define the submit handler
  const onSubmit: SubmitHandler<RegisterFormData> = async (data) => {
    // This is where you'll call your backend API for registration
    console.log("Registering user with data:", data);
    // Simulate network delay for loading state
    await new Promise(resolve => setTimeout(resolve, 2000));
    console.log("Pretend registration successful!");
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-900">Create a new account</h2>
        
        {/* 5. Use the handleSubmit wrapper */}
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            id="name"
            type="text"
            label="Full Name"
            register={register("name")}
            error={errors.name}
            autoComplete="name"
          />

          <Input
            id="email"
            type="email"
            label="Email Address"
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
            autoComplete="new-password"
          />

          <Input
            id="confirmPassword"
            type="password"
            label="Confirm Password"
            register={register("confirmPassword")}
            error={errors.confirmPassword}
            autoComplete="new-password"
          />

          <div>
            <Button type="submit" className="w-full" isLoading={isSubmitting}>
              Create Account
            </Button>
          </div>
        </form>

        <p className="text-sm text-center text-gray-600">
          Already have an account?{' '}
          <a href="/login" className="font-medium text-blue-600 hover:text-blue-500">
            Sign in
          </a>
        </p>
      </div>
    </div>
  );
}