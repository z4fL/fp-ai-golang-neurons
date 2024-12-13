import { useEffect } from "react";
import { useNavigate } from "react-router";

const AuthMiddleware = ({ children }) => {
  const navigate = useNavigate();
  const isAuthenticated = !!localStorage.getItem("session_token");

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
    }
  }, [isAuthenticated, navigate]);

  return isAuthenticated ? children : null;
};

export default AuthMiddleware;
