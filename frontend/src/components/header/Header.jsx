import Logout from "../auth/Logout";
import Navbar from "./Navbar";
import ToogleDarkMode from "./ToogleDarkMode";

const Header = () => {
  return (
    <header className="bg-gray-200 dark:bg-gray-900 p-4 flex justify-between items-center relative">
      <div className="absolute left-4">
        <Navbar />
      </div>
      <div className="mx-auto font-bold text-gray-800 dark:text-gray-200">
        Chatbot Smart Home Energy Management
      </div>
      <div className="absolute flex right-4">
        <ToogleDarkMode />
        <Logout />
      </div>
    </header>
  );
};

export default Header;
