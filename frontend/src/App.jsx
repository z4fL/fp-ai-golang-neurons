import React from "react";
import { useState, useEffect } from "react";
import ModalUpload from "./components/ModalUpload";

const App = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [request, setRequest] = useState("");
  const [chat, setChat] = useState([]);

  useEffect(() => {
    setChat([
      {
        id: 1,
        role: "assistant",
        content: "Halo, ada yang bisa aku bantu?",
        type: "text",
      },
      {
        id: 2,
        role: "user",
        content: "Iya, aku butuh bantuan buat layout.",
        type: "text",
      },
      {
        id: 3,
        role: "assistant",
        content:
          "Layout adalah cara di mana elemen-elemen disusun di dalam sebuah halaman web atau aplikasi. Ini mencakup penempatan, ukuran, dan tampilan elemen-elemen tersebut untuk menciptakan antarmuka pengguna yang intuitif dan menarik.",
        type: "text",
      }
    ]);
  }, []);

  const handleUploadFile = async () => {
    setChat((prevChat) => [
      ...prevChat,
      {
        id: prevChat.length + 1,
        role: "user",
        type: "file",
        name: file.name,
        size: file.size,
      },
    ]);
  };

  const handleChat = async () => {};

  return (
    <>
      <div className="flex flex-col h-screen bg-slate-50 font-noto">
        {/* Header */}
        <header className="bg-slate-800 text-white p-4 text-center font-bold">
          Chatbot Smart Home Energy Management
        </header>

        {/* Chat Area */}
        <main
          className="flex-1 overflow-y-auto"
          style={{ scrollbarGutter: "stable both-edges" }}
        >
          <div className="py-4">
            <div className="flex justify-center">
              <div
                id="chat-list"
                className="chat w-full max-w-screen-md flex flex-col space-y-4"
              >
                {chat.map((value) => (
                  <div
                    key={value.id}
                    className={`p-3 rounded-md ${
                      value.role === "assistant"
                        ? "self-start bg-slate-100 text-slate-900"
                        : "self-end bg-gray-800 text-white"
                    }`}
                  >
                    {value.role === "assistant" ? (
                      <p>{value.content}</p>
                    ) : value.type === "text" ? (
                      <p>{value.content}</p>
                    ) : (
                      <div>
                        <p>{value.name}</p>
                        <p className="text-sm text-slate-200">
                          {(value.size / 1024).toFixed(2)} KB
                        </p>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </main>

        {/* Input Area */}
        <footer className="w-full max-w-screen-md p-4 mb-4 mx-auto bg-gray-200 rounded-lg flex justify-center">
          <div className="w-full flex items-center space-x-2">
            {/* Upload File Button */}
            <button
              className="p-2 bg-lime-200 rounded-md hover:bg-lime-300"
              title="Upload File"
              onClick={() => setIsModalOpen(true)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="w-5 h-5 text-slate-900"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3 16.5V19a2 2 0 002 2h14a2 2 0 002-2v-2.5M16 9l-4-4m0 0l-4 4m4-4v12"
                />
              </svg>
            </button>

            {/* Input Field */}
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Ketik pesan..."
              className="flex-1 p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-lime-500 placeholder:text-slate-500"
            />

            {/* Send Button */}
            <button
              className="bg-lime-400 px-3 py-2 rounded-md hover:bg-lime-500"
              onClick={handleChat}
            >
              <svg
                className="w-6 h-6 text-slate-900"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 512 512"
              >
                <path
                  fill="currentColor"
                  d="M498.1 5.6c10.1 7 15.4 19.1 13.5 31.2l-64 416c-1.5 9.7-7.4 18.2-16 23s-18.9 5.4-28 1.6L284 427.7l-68.5 74.1c-8.9 9.7-22.9 12.9-35.2 8.1S160 493.2 160 480l0-83.6c0-4 1.5-7.8 4.2-10.8L331.8 202.8c5.8-6.3 5.6-16-.4-22s-15.7-6.4-22-.7L106 360.8 17.7 316.6C7.1 311.3 .3 300.7 0 288.9s5.9-22.8 16.1-28.7l448-256c10.7-6.1 23.9-5.5 34 1.4z"
                />
              </svg>
            </button>
          </div>
        </footer>
      </div>
      <ModalUpload
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        file={file}
        setFile={setFile}
        handleUpload={handleUploadFile}
      />
    </>
  );
};

export default App;
