import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";

import "./index.css";
import Login from "./components/Auth/Login.jsx";
import Register from "./components/Auth/Register.jsx";
import App from "./App.jsx";
import ChatArea from "./ChatArea.jsx";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<App />}>
          <Route path="/chats/:chatId" element={<ChatArea />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>
);
