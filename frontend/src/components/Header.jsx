import React from "react";
import Logout from "./Auth/Logout";

const Header = () => {
  return (
    <header className="bg-slate-800 text-white p-4 flex justify-center items-center relative">
      <div className="font-bold">Chatbot Smart Home Energy Management</div>
      <div className="absolute right-4">
        <Logout />
      </div>
    </header>
  );
};

export default Header;
