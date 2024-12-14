import { useNavigate } from "react-router";

const Logout = () => {
  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;
  const navigate = useNavigate();

  const handleLogout = async () => {
    const token = localStorage.getItem("session_token");
    if (token) {
      try {
        const response = await fetch(`${golangBaseUrl}/logout`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (response.ok) {
          localStorage.removeItem("session_token");
          navigate("/login");
        } else {
          console.error("Logout failed");
        }
      } catch (error) {
        console.error("Error:", error);
      }
    }
  };

  return <button onClick={handleLogout}>Logout</button>;
};

export default Logout;
