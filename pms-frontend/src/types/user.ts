// src/types/user.ts
export interface User {
  id: string; // UUID from backend
  hotel_id: string; // UUID from backend
  email: string;
  first_name?: string | null;
  last_name?: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  roles?: string[]; // e.g., ['admin', 'staff']
}