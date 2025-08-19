// src/components/ui/Input.tsx

import React from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";

// Define the props, including error and register from React Hook Form
interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  id: string;
  error?: FieldError;
  register: UseFormRegisterReturn; // Type for the register function
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, id, type = "text", error, register, ...props }, ref) => {
    const errorClass = error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500';

    return (
      <div>
        <label htmlFor={id} className="block text-sm font-medium text-gray-700">
          {label}
        </label>
        <div className="mt-1">
          <input
            {...props}
            {...register} // Spread the register props AFTER props to avoid override
            id={id}
            type={type}
            className={`block w-full px-3 py-2 border rounded-md shadow-sm sm:text-sm transition-colors duration-200 ${errorClass}`}
            aria-invalid={error ? 'true' : 'false'}
            aria-describedby={error ? `${id}-error` : undefined}
            ref={ref}     // Forward the ref
          />
        </div>
        {error && (
          <p
            id={`${id}-error`}
            className="mt-2 text-sm text-red-600"
            role="alert"
            aria-live="polite"
          >
            {error.message}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';