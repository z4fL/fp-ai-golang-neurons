import { useEffect, useState } from "react";
import fetchWithToken from "../../utility/fetchWithToken";
import { Link, useLocation } from "react-router";
import Bars from "../svg/Bars";
import SquarePlus from "../svg/SquarePlus";
import X from "../svg/X";

const Navbar = () => {
  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;
  const token = localStorage.getItem("session_token");
  const [listChats, setListChats] = useState([]);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const location = useLocation();

  useEffect(() => {
    if (isSidebarOpen) {
      setIsSidebarOpen(false);
    }

    const fetchListChats = async () => {
      try {
        const response = await fetchWithToken(
          `${golangBaseUrl}/chats`,
          undefined,
          token
        );

        if (!response.ok) throw new Error("failed to fetch list chat of user");

        const data = await response.json();
        setListChats(data.answer);
      } catch (error) {
        setListChats([]);
      }
    };

    fetchListChats();
  }, [location]);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <>
      <button className="py-2" onClick={toggleSidebar}>
        <Bars />
      </button>

      <button className="py-2 pl-6">
        <Link to="/">
          <SquarePlus />
        </Link>
      </button>

      {/* Sidebar */}
      <div
        className={`fixed top-0 left-0 z-10 h-full w-64 bg-gray-300 dark:bg-gray-950 text-gray-900 dark:text-gray-300 transform ${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        } transition-transform duration-300`}
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center space-x-4 p-2">
            <button className="p-2 rounded" onClick={toggleSidebar}>
              <X />
            </button>
            <h2 className="font-bold text-lg">List Chats</h2>
          </div>
          <div className="">
            <ul className="flex flex-col divide-y">
              {listChats.length ? (
                listChats.map((chat) => (
                  <Link key={chat.chatID} to={`/chats/${chat.chatID}`}>
                    <li className="px-4 py-2 hover:bg-lime-200 dark:hover:bg-lime-300 dark:hover:text-gray-800">
                      <p className="cursor-pointer truncate">{chat.content}</p>
                    </li>
                  </Link>
                ))
              ) : (
                <li className="px-4 py-2 dark:text-gray-300">
                  <p className="cursor-pointer truncate">No Chat</p>
                </li>
              )}
            </ul>
          </div>
        </div>
      </div>
    </>
  );
};

export default Navbar;
