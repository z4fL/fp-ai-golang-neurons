import { useEffect, useState } from "react";
import fetchWithToken from "../utility/fetchWithToken";
import { Link } from "react-router";

const Navbar = () => {
  const [listChats, setListChats] = useState([
    {
      chatID: 1,
      content: "From the provided data,...",
    },
    {
      chatID: 2,
      content: "From the provided data,...",
    },
  ]);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const token = localStorage.getItem("session_token");

  // useEffect(() => {
  //   const fetchListChats = async () => {
  //     try {
  //       const response = await fetchWithToken(
  //         "http://localhost:8080/chats",
  //         undefined,
  //         token
  //       );
  //       if (!response.ok) throw new Error("failed to fetch list chat of user");
  //       const data = await response.json();
  //       setListChats(data.answer);
  //     } catch (error) {
  //       console.error(error);
  //     }
  //   };

  //   fetchListChats();
  // }, []);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <>
      <button
        className="text-gray-800 font-bold py-2 px-4 rounded"
        onClick={toggleSidebar}
      >
        Navbar
      </button>
      {/* Sidebar */}
      <div
        className={`fixed top-0 left-0 z-10 h-full w-64 bg-gray-300 text-gray-800 transform ${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        } transition-transform duration-300`}
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center space-x-4 p-2">
            <button className="p-2 rounded" onClick={toggleSidebar}>
              Close
            </button>
            <h2 className="font-bold text-lg">List Chats</h2>
          </div>
          <div className="">
            <ul className="flex flex-col divide-y">
              {listChats.map((chat) => (
                <li key={chat.chatID} className="px-4 py-2 hover:bg-lime-200">
                  <Link to={`/chats/${chat.chatID}`}>
                    <div className="cursor-pointer">{chat.content}</div>
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </>
  );
};

export default Navbar;
