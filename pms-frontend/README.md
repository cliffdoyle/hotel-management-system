# Hotel Management System - Frontend

A modern React frontend for the Hotel Management System built with TypeScript, Vite, and Tailwind CSS.

## 🚀 Tech Stack

- **React 18** - Modern React with hooks and functional components
- **TypeScript** - Type safety and better developer experience
- **Vite** - Fast build tool and development server
- **Tailwind CSS** - Utility-first CSS framework
- **React Router v6** - Client-side routing
- **React Hook Form** - Form handling with validation
- **Zod** - Schema validation
- **Axios** - HTTP client for API calls

## 📁 Project Structure

```
src/
├── api/                 # API service layer
│   ├── client.ts       # Axios configuration
│   └── authService.ts  # Authentication API calls
├── components/         # Reusable UI components
│   └── ui/            # Basic UI components
├── context/           # React Context providers
├── hooks/             # Custom React hooks
├── pages/             # Page components
├── router/            # Route protection components
├── services/          # Service layer exports
├── types/             # TypeScript type definitions
└── main.tsx          # Application entry point
```

## 🛠️ Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env
   ```

   Update the `.env` file with your API endpoint:
   ```env
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   VITE_NODE_ENV=development
   ```

3. **Start Development Server**
   ```bash
   npm run dev
   ```

4. **Build for Production**
   ```bash
   npm run build
   ```

## 🔐 Authentication

The app includes a complete authentication system:

- **Login Page** - User authentication with form validation
- **Register Page** - User registration with comprehensive validation
- **Forgot Password** - Password reset functionality
- **Change Password** - Secure password change for authenticated users
- **Protected Routes** - Route protection based on authentication status
- **Token Management** - Automatic token handling with localStorage

## 🎯 Features Implemented

### ✅ Issue 33: Setup React Project Structure & Dependencies
- [x] React project with Vite
- [x] TypeScript configuration
- [x] Tailwind CSS setup
- [x] React Router installation
- [x] Axios for API communication
- [x] Folder structure (components, pages, hooks, services, types)
- [x] Environment variables configuration
- [x] ESLint and Prettier setup

### ✅ Issue 34: Authentication UI Components
- [x] Login component with form validation
- [x] Register component with comprehensive form fields
- [x] Password Change component
- [x] Forgot Password component
- [x] Form validation with proper error messages
- [x] Loading states and success/error feedback
- [x] Responsive design for mobile and desktop
- [x] Accessibility features (ARIA labels, keyboard navigation)

## 🚀 Next Steps

The frontend is ready for:
- Issue 35: Authentication Service & State Management (partially complete)
- Issue 36: Dashboard Layout & Navigation
- Issue 37: User Profile Management UI
- Additional hotel management features

## 🤝 Contributing

1. Follow the established folder structure
2. Use TypeScript for all new components
3. Implement proper error handling
4. Add accessibility features
5. Write responsive designs
6. Follow the existing code style

## 📄 License

This project is part of the Hotel Management System.
