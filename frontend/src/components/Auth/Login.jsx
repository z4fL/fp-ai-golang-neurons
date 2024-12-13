import React, { useState } from "react";
import { useNavigate } from "react-router";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isError, setIsError] = useState(false);

  const navigate = useNavigate();


  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const req = await fetch(`${golangBaseUrl}/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          password,
        }),
      });

      const res = await req.json()
      console.log(res.answer);
      
      if (!req.ok) throw new Error(res.answer);
      navigate("/");
    } catch (error) {
      setIsError(true);
      setUsername("")
      setPassword("")
      console.log(error);
    }
  };

  return (
    <div className="flex justify-center items-start md:items-center min-h-screen">
      <div className="flex max-w-screen-md w-full flex-col-reverse space-y-10 space-y-reverse md:space-y-0 md:flex-row md:justify-center mt-7 md:mb-20">
        <div className="md:w-1/2 flex justify-center font-noto">
          <div className="w-3/4 sm:w-1/2 md:w-3/4">
            <h2 className="text-2xl text-center font-bold mb-4">Login</h2>
            {isError && (
              <div className="mb-4 text-red-500 text-center">
                Invalid username or password
              </div>
            )}
            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <label
                  className="block text-sm font-medium mb-2"
                  htmlFor="username"
                >
                  Username
                </label>
                <input
                  className="w-full px-3 py-2 border rounded"
                  type="text"
                  id="username"
                  name="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
              <div className="mb-4">
                <label
                  className="block text-sm font-medium mb-2"
                  htmlFor="password"
                >
                  Password
                </label>
                <input
                  className="w-full px-3 py-2 border rounded"
                  type="password"
                  id="password"
                  name="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <button
                className="w-full bg-lime-500 hover:bg-lime-600 active:bg-lime-700 disabled:bg-lime-500 text-white py-2 rounded"
                type="submit"
                disabled={!username.trim() || !password.trim()}
              >
                Submit
              </button>
            </form>
            <div className="text-center mt-4">
              <span>Not have an account? </span>
              <button
                className="text-lime-800 hover:underline"
                onClick={() => navigate('/register')}
              >
                Register
              </button>
            </div>
          </div>
        </div>
        <div className="md:w-1/2 flex justify-center items-center">
          <img
            src="/logo.png"
            alt="Logo"
            className="max-w-20 md:max-w-full h-auto"
          />
        </div>
      </div>
    </div>
  );
};

export default Login;
