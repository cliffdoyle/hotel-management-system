# File: backend/README.md

# Hotel Management System - Backend

This is the backend for the Hotel Management System, a MEWS-like SaaS application. It is built with Go, using a clean architecture and standard library conventions.

## Tech Stack
- **Language**: Go
- **Router**: httprouter
- **Database**: Supabase (PostgreSQL)
- **Caching**: Redis

## Setup

1.  **Clone the repository.**
2.  **Navigate to the `backend` directory:**
    ```sh
    cd backend
    ```
3.  **Create an environment file:**
    ```sh
    cp .env.example .env
    ```
4.  **Update `.env`:**
    Fill in your Supabase Database DSN and other variables.
5.  **Install Dependencies:**
    ```sh
    go mod tidy
    ```
6.  **Run the application:**
    ```sh
    go run ./cmd/api
    ```
The server will start on the port specified in your `.env` file (default is 8080). You can check its status at `http://localhost:8080/v1/healthcheck`.