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
    const savedChat = localStorage.getItem("chatSession");

    const initChat = {
      id: 1,
      role: "assistant",
      content: "Hello, how can I help you?",
    };

    setIsReload(true);

    if (savedChat) {
      return JSON.parse(savedChat);
    } else {
      return [initChat];
    }
  };

  const [chatHistory, setChatHistory] = useState(initializeChatHistory);

  const resetChatSession = async () => {
    localStorage.removeItem("chatSession");
    localStorage.removeItem("lastAccess");
    const res = await fetch(`${golangBaseUrl}/remove-session`, {
      method: "POST",
    });
    if (!res.ok) console.log("Failed to remove session");
    else console.log("success to remove session");
    setChatHistory(initializeChatHistory());
  };

  useEffect(() => {
    const interval = setInterval(() => {
      const lastAccess = localStorage.getItem("lastAccess");
      console.log("Checking session timeout...");

      if (
        lastAccess &&
        Date.now() - parseInt(lastAccess, 10) > 30 * 60 * 1000
      ) {
        console.log("Session expired. Resetting chat session.");
        resetChatSession();
        clearInterval(interval);
      }
    }, 30 * 1000);

    return () => clearInterval(interval);
  }, [resetChatSession]);

  const appendChat2Session = (newChat) => {
    localStorage.setItem("chatSession", JSON.stringify(newChat));
    localStorage.setItem("lastAccess", Date.now());
  };

  const handleResponse = async () => {
    setIsLoading(true);
    const lastChat = chatHistory[chatHistory.length - 1];

    try {
      const res =
        lastChat.type === "text"
          ? await handleChat()
          : await handleUploadFile();

      const data = await res.json();

      if (!res.ok) throw new Error("Failed to fetch response");

      const responseChat = {
        id: chatHistory.length + 1,
        role: "assistant",
        content: data.answer,
        type: "text",
      };

      // remove LOADING... chat and add responseChat
      setChatHistory((prevChat) => {
        const updatedChat = [...prevChat.slice(0, -1), responseChat];
        appendChat2Session(updatedChat);
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
    const lastChat = chatHistory.at(-1);

    if (chatHistory.length && lastChat?.role === "user") {
      setChatHistory((prevChat) => [
        ...prevChat,
        {
          id: prevChat.length + 1,
          role: "assistant",
          content: "LOADING...",
          type: "text",
        },
      ]);

      handleResponse();
    }
  }, [chatHistory]);

  const handleChat = () => {
    const lastChat = chatHistory[chatHistory.length - 1];
    const previousChat = chatHistory[chatHistory.length - 2];

    const payload = {
      type: lastChat.content.includes("/file") ? "tapas" : "phi",
      query: lastChat.content.replace("/file", "").trim(),
      ...(previousChat.id !== 1 && { prevChat: previousChat.content }),
    };

    return fetch(`${golangBaseUrl}/chat`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
  };

  const handleUploadFile = () => {
    const formData = new FormData();
    formData.append("file", file);

    return fetch(`${golangBaseUrl}/upload`, {
      method: "POST",
      body: formData,
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
