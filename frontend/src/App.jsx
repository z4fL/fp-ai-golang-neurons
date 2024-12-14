import React, { useState, useEffect } from "react";
import ModalUpload from "./components/ModalUpload";
import ChatList from "./components/ChatList";
import Header from "./components/Header";
import Footer from "./components/Footer";

const App = () => {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const [isReload, setIsReload] = useState(false);

  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;

  const initializeChatHistory = () => {
    const savedChat = localStorage.getItem("chat_session");

    const initChat = {
      id: 1,
      role: "assistant",
      content: "Hello, how can I help you?",
    };

    if (savedChat) {
      setIsReload(true);
      return JSON.parse(savedChat);
    } else {
      return [initChat];
    }
  };

  const [chatHistory, setChatHistory] = useState(initializeChatHistory);

  const handleResponse = async () => {
    setIsLoading(true);
    try {
      const responseChat = await fetchChatResponse();

      // remove LOADING... chat and add responseChat
      setChatHistory((prevChat) => {
        const updatedChat = [...prevChat.slice(0, -1), responseChat];
        localStorage.setItem("chat_session", JSON.stringify(updatedChat));
        return updatedChat;
      });

      if (!file) setFile(null); // remove file
      setIsError(false);
    } catch (error) {
      // remove LOADING... chat and error chat
      setChatHistory((prevChat) => [
        ...prevChat.slice(0, -1),
        {
          id: prevChat.length + 1,
          role: "assistant",
          content: String(error),
          type: "error",
        },
      ]);

      setIsError(true);
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const lastChat = chatHistory.at(-1); // last element
    if (chatHistory.length && lastChat.role === "user") {
      setChatHistory((prevChat) => {
        return [
          ...prevChat,
          {
            id: prevChat.length + 1,
            role: "assistant",
            content: "LOADING...",
            type: "text",
          },
        ];
      });

      handleResponse();
    }
  }, [chatHistory]);

  const fetchChatResponse = async () => {
    const lastChat = chatHistory[chatHistory.length - 1];
    const res =
      lastChat.type === "text" ? await handleChat() : await handleUploadFile();

    const data = await res.json();

    if (!res.ok) throw new Error("Failed to fetch response");

    return {
      id: chatHistory.length + 1,
      role: "assistant",
      content: data.answer,
      type: "text",
    };
  };

  const handleChat = () => {
    const lastChat = chatHistory[chatHistory.length - 1];
    const previousChat = chatHistory[chatHistory.length - 2];

    const payload = {
      type: lastChat.content.includes("/file") ? "tapas" : "phi",
      query: lastChat.content.replace("/file", "").trim(),
      ...(previousChat.id !== 1 && { prevChat: previousChat.content }),
    };

    console.log("Payload for chat:", payload);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          json: () =>
            Promise.resolve({ answer: "Dummy fetch response for chat" }),
          ok: true,
        });
      }, 5000);
    });
  };

  const handleUploadFile = () => {
    const formData = new FormData();
    formData.append("file", file);

    console.log("FormData for file upload:", formData);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          json: () =>
            Promise.resolve({ answer: "Dummy fetch response for file upload" }),
          ok: true,
        });
      }, 5000);
    });
  };

  const getResponse = (type) => {
    setIsReload(false);
    if (type === "file" && !file) return;

    const newChat = {
      id: chatHistory.length + 1,
      role: "user",
      type,
      content: type === "text" ? query : { name: file.name, size: file.size },
    };

    if (type === "text") setQuery("");
    setIsLoading(true);
    setChatHistory((prevchat) => [...prevchat, newChat]);
  };

  const reloadChat = () => {
    if (chatHistory[chatHistory.length - 1].type === "error") {
      setChatHistory((prevChat) => prevChat.slice(0, -1));
    }
  };

  return (
    <>
      <div className="flex flex-col h-screen bg-slate-50 font-noto">
        <Header />

        <ChatList
          chatList={chatHistory}
          setIsLoading={setIsLoading}
          reloadChat={reloadChat}
          isReload={isReload}
        />

        {/* Input Area */}
        <Footer
          setIsModalOpen={setIsModalOpen}
          query={query}
          setQuery={setQuery}
          getResponse={getResponse}
          isLoading={isLoading}
          isError={isError}
        />
      </div>
      <ModalUpload
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        file={file}
        setFile={setFile}
        getResponse={() => getResponse("file")}
      />
    </>
  );
};

export default App;
