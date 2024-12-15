import Logout from "../auth/Logout";
import Navbar from "./Navbar";

const Header = () => {
  return (
    <header className="bg-gray-200 text-white p-4 flex justify-between items-center relative">
      <div className="absolute left-4">
        <Navbar />
      </div>
      <div className="mx-auto font-bold text-gray-800">
        Chatbot Smart Home Energy Management
      </div>
      <div className="absolute right-4">
        <Logout />
      </div>
    </header>
  );
};

export default Header;
