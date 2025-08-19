// src/api/userService.ts
import apiClient from "./client";
import type { User } from "../types/user";

interface UserProfileResponse {
  user: User;
}

export const getUserProfile = async (): Promise<User> => {
  try {
    const response = await apiClient.get<UserProfileResponse>("/users/profile");
    return response.data.user;
  } catch (error: unknown) {
    if (error && typeof error === 'object' && 'response' in error) {
      const axiosError = error as { 
        response?: { 
          data?: { 
            message?: string;
            error?: string;
          } 
        } 
      };
      
      const errorData = axiosError.response?.data;
      throw new Error(errorData?.message || errorData?.error || "Failed to fetch user profile.");
    }
    throw new Error("Failed to fetch user profile.");
  }
};
