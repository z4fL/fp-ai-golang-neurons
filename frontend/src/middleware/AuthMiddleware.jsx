import { useEffect } from "react";
import { useNavigate } from "react-router";
import Cookies from "js-cookie";

const AuthMiddleware = ({ children }) => {
  const navigate = useNavigate();
  const isAuthenticated = !!Cookies.get("session_token");

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
    }
  }, [isAuthenticated, navigate]);

  return isAuthenticated ? children : null;
};

export default AuthMiddleware;
