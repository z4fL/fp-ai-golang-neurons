import { useState, useEffect } from "react";
import { Outlet, useParams } from "react-router";

import Navbar from "./components/Navbar";
import Logout from "./components/Auth/Logout";
import SendSVG from "./components/svg/SendSVG";
import UploadSVG from "./components/svg/UploadSVG";
import ModalUpload from "./components/ModalUpload";
import Markdown from "react-markdown";

const App = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [file, setFile] = useState(null);
  const [isNewChat, setIsNewChat] = useState(false);

  const initChat = {
    id: 1,
    role: "assistant",
    content: "Hello, how can I help you?",
  };
  
  const { chatId } = useParams();

  useEffect(() => {
    if (chatId) {
      setIsNewChat(true);
    }
  }, []);

  return (
    <div className="relative flex flex-col h-screen bg-gray-50 font-noto">
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

      <main
        className="flex-1 overflow-y-auto"
        style={{ scrollbarGutter: "stable both-edges" }}
      >
        <div className="py-4">
          <div className="flex justify-center">
            {isNewChat ? (
              <Outlet />
            ) : (
              <div
                id="chat-list"
                className="chat w-full max-w-screen-md flex flex-col space-y-4"
              >
                <div
                  className={`self-start p-3 rounded-md bg-slate-100 text-slate-900`}
                >
                  <div className="prose prose-base">
                    <Markdown>{initChat.content}</Markdown>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>

      <footer className="w-full max-w-screen-md p-4 mb-4 mx-auto bg-gray-200 rounded-lg flex justify-center">
        <div className="w-full flex items-center space-x-2">
          {/* Upload File Button */}
          <button
            className="p-2 bg-lime-200 rounded-md hover:bg-lime-300"
            title="Upload File"
          >
            <UploadSVG />
          </button>

          {/* Input Field */}
          <input
            type="text"
            placeholder="Ketik pesan..."
            className="flex-1 p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-lime-500 placeholder:text-slate-500"
          />

          {/* Send Button */}
          <button
            className={`bg-lime-400 px-3 py-2 rounded-md hover:bg-lime-500`}
          >
            <SendSVG />
          </button>
        </div>
      </footer>
      <ModalUpload
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        file={file}
        setFile={setFile}
        getResponse={() => getResponse("file")}
      />
    </div>
  );
};

export default App;
