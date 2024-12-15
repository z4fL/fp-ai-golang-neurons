import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";

import "./index.css";
import Login from "./components/auth/Login.jsx";
import Register from "./components/auth/Register.jsx";
import ChatArea from "./components/ChatArea.jsx";
import ChatLayout from "./Layout/ChatLayout.jsx";
import NewChat from "./components/NewChat.jsx";
import AuthMiddleware from "./middleware/AuthMiddleware.jsx";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route
          element={
            <AuthMiddleware>
              <ChatLayout />
            </AuthMiddleware>
          }
        >
          <Route exact path="/" element={<NewChat />} />
          <Route path="/chats/:chatId" element={<ChatArea />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>
);
