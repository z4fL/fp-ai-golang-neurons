import React, { useState, useEffect } from "react";
import ModalUpload from "./components/ModalUpload";
import ChatList from "./components/ChatList";
import Header from "./components/Header";
import Footer from "./components/Footer";
import fetchWithToken from "./utility/fetchWithToken";
import LoadChat from "./components/LoadChat";

const AppBackup = () => {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const [isReload, setIsReload] = useState(false);
  const [isGetChatHistory, setIsGetChatHistory] = useState(false);

  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;

  const token = localStorage.getItem("session_token");

  const initChat = {
    id: 1,
    role: "assistant",
    content: "Hello, how can I help you?",
  };

  const [chatHistory, setChatHistory] = useState([initChat]);

  const initializeChatHistory = async () => {
    setIsGetChatHistory(true);

    try {
      const res = await fetchWithToken(
        `${golangBaseUrl}/chats`,
        { method: "GET" },
        token
      );

      const { answer } = await res.json();
      if (!res.ok) throw new Error(answer);
      const savedChat = answer;
      console.log(savedChat);

      if (savedChat && savedChat.length > 0) {
        setIsReload(true);
        return savedChat;
      } else {
        return [initChat];
      }
    } catch (error) {
      return [initChat];
    } finally {
      setIsGetChatHistory(false);
    }
  };

  useEffect(() => {
    const fetchChatHistory = async () => {
      const history = await initializeChatHistory();
      setChatHistory(history);
    };

    fetchChatHistory();
  }, []);

  const handleResponse = async () => {
    setIsLoading(true);
    try {
      const responseChat = await fetchChatResponse();

      // remove LOADING... chat and add responseChat
      setChatHistory((prevChat) => [...prevChat.slice(0, -1), responseChat]);

      if (chatHistory.length <= 3) {
        await createChatHandler(responseChat);
      } else {
        await addMessageHandler(responseChat);
      }

      if (!file) setFile(null); // remove file
      setIsError(false);
    } catch (error) {
      const errorChat = {
        id: chatHistory.length + 1,
        role: "assistant",
        content: String(error),
        type: "error",
      };

      // remove LOADING... chat and error chat
      setChatHistory((prevChat) => [...prevChat.slice(0, -1), errorChat]);

      if (chatHistory.length <= 3) {
        await createChatHandler(errorChat);
      } else {
        await addMessageHandler(errorChat);
      }

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

    if (!res.ok) throw new Error(data.answer);

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

    return fetchWithToken(
      `${golangBaseUrl}/chat-with-ai`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      },
      token
    );
  };

  const handleUploadFile = () => {
    const formData = new FormData();
    formData.append("file", file);

    return fetchWithToken(
      `${golangBaseUrl}/upload`,
      {
        method: "POST",
        body: formData,
      },
      token
    );
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

  // Fungsi untuk memanggil handler CreateChat
  const createChatHandler = async (responseChat) => {
    const payload = {
      chat_history: [...chatHistory, responseChat], // Kirim chat history yang sudah ada
    };

    console.log(payload);

    const res = await fetchWithToken(
      `${golangBaseUrl}/chats`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      },
      token
    );

    if (!res.ok) {
      console.log("Failed to create chat");
    } else {
      console.log("Chat created successfully");
    }
  };

  // Fungsi untuk memanggil handler AddMessage
  const addMessageHandler = async (responseChat) => {
    const lastChat = chatHistory.at(-1);
    const payload = {
      chat_history: [lastChat, responseChat], // Kirim chat history yang sudah ada
    };
    console.log(payload);

    const res = await fetchWithToken(
      `${golangBaseUrl}/chats`,
      {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      },
      token
    );

    if (!res.ok) {
      console.log("Failed to add message");
    } else {
      console.log("Message added successfully");
    }
  };

  return (
    <>
      <div className="flex flex-col h-screen bg-slate-50 font-noto">
        <Header />

        {isGetChatHistory ? (
          <LoadChat />
        ) : (
          <ChatList
            setIsError={setIsError}
            chatList={chatHistory}
            setIsLoading={setIsLoading}
            reloadChat={reloadChat}
            isReload={isReload}
          />
        )}

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

export default AppBackup;
