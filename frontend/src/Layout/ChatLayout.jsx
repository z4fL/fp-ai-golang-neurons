import { useState, useEffect } from "react";
import { Outlet, useNavigate, useParams } from "react-router";
import Header from "../components/header/Header";
import Footer from "../components/footer/Footer";
import ModalUpload from "../components/ModalUpload";
import fetchWithToken from "../utility/fetchWithToken";
import NewChat from "../components/NewChat";

const ChatLayout = () => {
  const golangBaseUrl = import.meta.env.VITE_GOLANG_URL;
  const token = localStorage.getItem("session_token");
  const navigate = useNavigate();
  const { chatId } = useParams();

  const initChat = {
    id: 1,
    role: "assistant",
    content: "Hello, how can I help you?",
    type: "text",
  };

  const [chatHistory, setChatHistory] = useState([initChat]);

  const [query, setQuery] = useState("");
  const [file, setFile] = useState(null);

  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [errorType, setErrorType] = useState("");

  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    if (chatId) {
      const fetchChat = async () => {
        try {
          const res = await fetchWithToken(
            `${golangBaseUrl}/chats/${chatId}`,
            undefined,
            token
          );

          if (!res.ok) {
            throw new Error("Failed to fetch chat");
          }

          const data = await res.json();
          setChatHistory(data.answer);

          setErrorType("");
          setIsLoading(false);
          setIsError(false);
          setQuery("");
          setFile(null);
        } catch (error) {
          console.error(error);
          setIsError(true);
        }
      };

      fetchChat();
    } else {
      setChatHistory([initChat]);
      setErrorType("");
      setIsLoading(false);
      setIsError(false);
      setQuery("");
      setFile(null);
    }
  }, [chatId]);

  const getResponse = (type) => {
    if (type === "file" && !file) return;

    const newChat = {
      id: chatHistory.length + 1,
      role: "user",
      type,
      content: type === "text" ? query : { name: file.name, size: file.size },
    };

    if (type === "text") setQuery("");
    setChatHistory((prevchat) => [...prevchat, newChat]);
  };

  const handleResponse = async () => {
    setIsLoading(() => {
      console.log("ChatLayout isLoading:", true);
      return true;
    });
    try {
      const responseChat = await fetchChatResponse();
      console.log("errorType :", errorType);
      
      // remove LOADING... chat and add responseChat
      setChatHistory((prevChat) => [...prevChat.slice(0, -1), responseChat]);

      if (chatHistory.length <= 3 && errorType === "") {
        console.log("createNewChat");

        await createNewChat(responseChat);
      } else {
        console.log("updateChat");
        await updateChat(responseChat);
      }

      setIsError(false);
    } catch (error) {
      console.log(error);

      const errorChat = {
        id: chatHistory.length + 1,
        role: "assistant",
        content: String(error),
        type: "error",
      };

      // remove LOADING... chat and error chat
      setChatHistory((prevChat) => [...prevChat.slice(0, -1), errorChat]);

      if (chatHistory.length <= 3) {
        if (!error) {
          console.log("createNewChat");
          await createNewChat(errorChat);
        }
      } else {
        console.log("updateChat");
        await updateChat(errorChat);
      }

      setIsError(true);
      setIsLoading(false);
    }
  };

  useEffect(() => {
    console.log(chatHistory.length);

    const lastChat = chatHistory.at(-1); // last element
    if (lastChat.role === "user") {
      setChatHistory((prevChat) => [
        ...prevChat,
        {
          id: prevChat.length + 1,
          role: "assistant",
          content: "LOADING...",
          type: "loading",
        },
      ]);

      handleResponse();
    }
  }, [chatHistory]);

  async function fetchChatResponse() {
    const lastChat = chatHistory[chatHistory.length - 1];
    const res =
      lastChat.type === "text" ? await handleChat() : await handleUploadFile();

    const data = await res.json();

    if (file) setFile(null); // remove file
    if (!res.ok) throw new Error(data.answer);

    return {
      id: chatHistory.length + 1,
      role: "assistant",
      content: data.answer,
      type: "text",
    };
  }

  async function handleChat() {
    const lastChat = chatHistory[chatHistory.length - 1];
    const previousChat = chatHistory[chatHistory.length - 2];

    const payload = {
      type: lastChat.content.includes("/file") ? "tapas" : "phi",
      query: lastChat.content.replace("/file", "").trim(),
      ...(previousChat.id !== 1 && { prevChat: previousChat.content }),
    };

    const res = await fetchWithToken(
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

    if (!res.ok) {
      setErrorType("text");
    }

    return res;
  }

  async function handleUploadFile() {
    const formData = new FormData();
    formData.append("file", file);

    const res = await fetchWithToken(
      `${golangBaseUrl}/upload`,
      {
        method: "POST",
        body: formData,
      },
      token
    );

    if (!res.ok) {
      setErrorType("file");
    }

    return res;
  }

  async function createNewChat(responseChat) {
    const payload = {
      chat_history: [...chatHistory, responseChat], // Kirim chat history yang sudah ada
    };

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
      console.log("Sucess to create chat");
      const data = await res.json();
      if (!file) setFile(null); // remove file
      setIsError(false);

      navigate(`/chats/${data.answer}`, { state: { fromNavigate: true } });
    }
  }

  async function updateChat(responseChat) {
    const lastChat = chatHistory.at(-1);
    const payload = {
      chat_history: [lastChat, responseChat], // Kirim chat history yang sudah ada
    };

    const res = await fetchWithToken(
      `${golangBaseUrl}/chats/${chatId}`,
      {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      },
      token
    );

    if (!res.ok) {
      console.log("Failed to add chat");
    } else {
      console.log("Chat added successfully");
    }
  }

  const reloadChat = () => {
    console.log("Reloading chat");
    if (chatHistory[chatHistory.length - 1].type === "error") {
      if (
        chatHistory.length > 1 &&
        chatHistory[chatHistory.length - 2].type === "file"
      ) {
        setChatHistory((prevChat) => prevChat.slice(0, -2));
        setIsError(false);
      } else {
        setChatHistory((prevChat) => prevChat.slice(0, -1));
      }
    }
  };

  return (
    <div className="relative flex flex-col h-screen bg-gray-50 dark:bg-gray-800 font-noto">
      <Header />
      {chatId ? (
        <Outlet
          context={{ chatHistory, setIsLoading, setIsError, reloadChat }}
        />
      ) : (
        <NewChat
          chatList={chatHistory}
          reloadChat={reloadChat}
          setIsError={setIsError}
          setIsLoading={setIsLoading}
        />
      )}

      <Footer
        setIsModalOpen={setIsModalOpen}
        query={query}
        setQuery={setQuery}
        getResponse={() => getResponse("text")}
        isLoading={isLoading}
        isError={isError}
        errorType={errorType}
      />
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

export default ChatLayout;
