import React, { useState, useEffect } from "react";
import ModalUpload from "./components/ModalUpload";
import ChatList from "./components/ChatList";
import Header from "./components/Header";
import Footer from "./components/Footer";

const App = () => {
  const [file, setFile] = useState(null); // file user
  const [query, setQuery] = useState(""); // query user

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const [chatHistory, setChatHistory] = useState([
    {
      id: 1,
      role: "assistant",
      content: "Halo, ada yang bisa aku bantu?",
    },
  ]);

  const handleResponse = async () => {
    setIsLoading(true);
    const lastChat = chatHistory[chatHistory.length - 1];

    try {
      const res =
        lastChat.type === "text"
          ? await handleChat()
          : await handleUploadFile();

      if (!res.ok) {
        setIsError(true);
        throw new Error("Failed to fetch response");
      }

      const data = await res.json();
      setChatHistory((prevChat) => prevChat.slice(0, -1));

      setChatHistory((prevChat) => [
        ...prevChat,
        {
          id: prevChat.length + 1,
          role: "assistant",
          content: data.answer,
          type: "text",
        },
      ]);

      setFile(null);
      setIsError(false);
    } catch (error) {
      setIsLoading(false);
      setChatHistory((prevChat) => prevChat.slice(0, -1));

      setChatHistory((prevChat) => [
        ...prevChat,
        {
          id: prevChat.length + 1,
          role: "assistant",
          content: "ERROR",
          type: "error",
        },
      ]);

      setIsError(true);
    }
  };

  useEffect(() => {
    if (
      chatHistory.length &&
      chatHistory[chatHistory.length - 1].role === "user"
    ) {
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

  const handleChat = async () => {
    const lastChat = chatHistory[chatHistory.length - 1];
    const previousChat = chatHistory[chatHistory.length - 2];

    const queryTapas = lastChat.content.includes("/file")
      ? lastChat.content.replace("/file", "").trim()
      : lastChat.content;

    const payload = {
      type: lastChat.content.includes("/file") ? "tapas" : "phi",
      query: queryTapas,
      ...(previousChat?.id !== 1 && { prevChat: previousChat?.content }),
    };

    return fetch("http://localhost:8080/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
  };

  const handleUploadFile = () => {
    const formData = new FormData();
    formData.append("file", file);

    return fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData,
    });
  };

  const getResponse = (type) => {
    if (type === "file") if (!file) return;

    const newChat = {
      id: chatHistory.length + 1,
      role: "user",
      type,
      content: type === "text" ? query : { name: file?.name, size: file?.size },
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
        {/* Header */}
        <Header />

        {/* Chat Area */}
        <ChatList
          chatList={chatHistory}
          setIsLoading={setIsLoading}
          reloadChat={reloadChat}
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
