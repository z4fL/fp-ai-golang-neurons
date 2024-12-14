import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

const AuthMiddleware = ({ children }) => {
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const checkToken = async () => {
      const token = localStorage.getItem("session_token");

      if (!token) {
        navigate("/login");
        return;
      }

      try {
        const response = await fetch("http://localhost:8080/validate-session", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          localStorage.removeItem("session_token");
          setIsLoading(false);
          setError("Session expired. Please log in again.");
          setTimeout(() => {
            navigate("/login");
          }, 3000);
        } else if (response.ok) {
          setIsLoading(false);
        } else {
          throw new Error("Unexpected error");
        }
      } catch (err) {
        console.error(err);
        setError("An error occurred. Please try again.");
      }
    };

    checkToken();
  }, []);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        Loading...
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        {error}
      </div>
    );
  }

  return children;
};

export default AuthMiddleware;
