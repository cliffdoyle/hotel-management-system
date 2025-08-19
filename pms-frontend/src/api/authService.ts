// src/api/authService.ts
import apiClient from "./client";
import type { User } from "../types/user";

// Define the shape of your form data
import type { LoginFormData } from "../pages/LoginPage";
import type { RegisterFormData } from "../pages/RegisterPage";

// Define the expected API response for auth endpoints
interface RegisterResponse {
  user: User;
}

interface LoginResponse {
  tokens: {
    authentication_token: {
      token: string;
      expiry: string;
    };
    refresh_token: string;
  };
}

// Unified response for frontend use
interface AuthResponse {
  token: string;
  user: User;
}

export const loginUser = async (credentials: LoginFormData): Promise<AuthResponse> => {
  try {
    const response = await apiClient.post<LoginResponse>("/users/login", credentials);

    // Transform backend response to frontend format
    const { tokens } = response.data;

    // We need to get user data from the token or make another API call
    // For now, we'll create a minimal user object - this should be improved
    // when we have a proper user profile endpoint
    const user: User = {
      id: '', // Will be populated from token validation
      hotel_id: '',
      email: credentials.email,
      first_name: null,
      last_name: null,
      is_active: true,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      roles: []
    };

    return {
      token: tokens.authentication_token.token,
      user
    };
  } catch (error: unknown) {
    if (error && typeof error === 'object' && 'response' in error) {
      const axiosError = error as {
        response?: {
          data?: {
            message?: string;
            error?: string;
            errors?: Record<string, string>;
          }
        }
      };

      const errorData = axiosError.response?.data;

      // Handle validation errors
      if (errorData?.errors) {
        const validationMessages = Object.values(errorData.errors).join(', ');
        throw new Error(validationMessages);
      }

      // Handle general error messages
      throw new Error(errorData?.message || errorData?.error || "Invalid email or password.");
    }
    throw new Error("Invalid email or password.");
  }
};

export const registerUser = async (data: RegisterFormData): Promise<AuthResponse> => {
  try {
    const response = await apiClient.post<RegisterResponse>("/users/register", data);

    // Backend returns user but no token for registration
    // We need to login after successful registration to get tokens
    const loginResponse = await loginUser({
      email: data.email,
      password: data.password
    });

    // Update user data with registration response
    return {
      token: loginResponse.token,
      user: {
        ...response.data.user,
        // Ensure we have the user data from registration
      }
    };
  } catch (error: unknown) {
    if (error && typeof error === 'object' && 'response' in error) {
      const axiosError = error as {
        response?: {
          data?: {
            message?: string;
            error?: string;
            errors?: Record<string, string>;
          }
        }
      };

      const errorData = axiosError.response?.data;

      // Handle validation errors
      if (errorData?.errors) {
        const validationMessages = Object.values(errorData.errors).join(', ');
        throw new Error(validationMessages);
      }

      // Handle general error messages
      throw new Error(errorData?.message || errorData?.error || "Registration failed. Please try again.");
    }
    throw new Error("Registration failed. Please try again.");
  }
};