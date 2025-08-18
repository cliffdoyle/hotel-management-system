// src/App.tsx

import { BrowserRouter, Routes, Route } from "react-router-dom";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { ForgotPasswordPage } from "./pages/ForgotPasswordPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        {/* We will add other routes here later, e.g., Register, Dashboard */} 
        <Route path="*" element={<LoginPage />} /> {/* Default to login */}
      </Routes>
    </BrowserRouter>
  );
}

export default App;