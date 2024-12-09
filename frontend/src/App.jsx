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
    const previousChat = chatHistory[chatHistory.length - 2] || "";

    let payload = {};

    let queryTapas = lastChat.content;

    // Cek jika chat terakhir berisi "/file"
    if (queryTapas.includes("/file")) {
      queryTapas = queryTapas.replace("/file", "").trim();
    }

    payload =
      chatHistory[chatHistory.length - 2].id === 1 || payload.type === "tapas"
        ? { type: "tapas", query: queryTapas }
        : {
            type: "phi",
            query: lastChat.content,
            prevChat: previousChat.content,
          };

    return new Promise((resolve) => {
      console.log("Request Payload (Chat):\n", payload);
      setTimeout(() => {
        resolve({
          ok: true,
          json: () => ({
            status: "success",
            answer:
              "The EVCar, or Electric Vehicle Charger, uses more electricity compared to the other devices because it is designed to draw a high amount of power for a specific purpose. Electric vehicles require substantial electrical charge to run their batteries, and charging these batteries consumes a significant amount of electricity.\n\nHere's a detailed breakdown of why it uses more electricity:\n\n1. Battery size: Electric cars typically have larger batteries than other household electronic devices. The battery is the main energy storage component in an EV, and its size directly corresponds to the electricity the vehicle will consume while charging. In essence, the larger the battery, the more electricity required to fully charge it.\n\n2. Charging Power: Different EVs have varying charging power requirements (measured in kW), and chargers need to match this specification to charge an EV efficiently. While some devices may consume low-power electricity, EV chargers require a more substantial power flow to charge the car's battery quickly.\n\n3. Charge Time: The time required to charge an electric vehicle greatly depends on the battery's capacity and the charging power. EV charging times can range from several hours to overnight (possibly up to 22 hours for some models). During this time, the EV charger continuously operates, consuming a steady flow of electricity.\n\n4. Energy Demand: Due to the concept of duty cycles in electronics — where devices operate at their peak capacity over longer periods — the constant operation of the EV charger signifies a higher energy demand compared to devices like mobile phones, laptops, or even TVs which may have periods of low or no usage.\n\nIn summary, the high energy consumption of EV chargers is a direct outcome of their purpose: they must supply a substantial amount of electricity over a sustained period to recharge electric vehicle batteries. This energy demand far exceeds that of other more conventional electronic devices used in a household setting.",
          }),
        });
      }, 5000);
    });
  };

  const handleUploadFile = () => {
    return new Promise((resolve) => {
      console.log("Request Payload (File):\n", {
        file,
      });
      setTimeout(() => {
        resolve({
          ok: true,
          json: () => ({
            status: "success",
            answer:
              "From the provided data, here are the Least Electricity: TV and the Most Electricity: EVCar.",
          }),
        });
      }, 5000);
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
