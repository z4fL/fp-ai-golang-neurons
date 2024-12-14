import React, { useState, useRef } from "react";
import { useEffect } from "react";
import UserChat from "./Chat/UserChat";
import AssistantChat from "./Chat/AssistantChat";

const ChatList = ({
  chatList,
  setIsLoading,
  reloadChat,
  isReload,
  setIsError,
}) => {
  const [displayResponse, setDisplayResponse] = useState("");
  const [isCompletedTyping, setIsCompletedTyping] = useState(false);

  const [isAutoScrollEnabled, setIsAutoScrollEnabled] = useState(true);
  const bottomRef = useRef(null);
  const chatContainerRef = useRef(null);

  useEffect(() => {
    const lastChat = chatList[chatList.length - 1];

    if (lastChat.role !== "assistant") return;
    if (isReload) {
      setDisplayResponse(lastChat.content);
      setIsCompletedTyping(true);
      return;
    }

    setIsCompletedTyping(false);

    let i = 0;
    const responseAssistant = lastChat.content;

    const intervalId = setInterval(() => {
      setDisplayResponse(responseAssistant.slice(0, i));
      i++;

      if (i > responseAssistant.length) {
        setIsLoading(false);

        clearInterval(intervalId);
        setIsCompletedTyping(true);
      }
    }, 5);

    return () => clearInterval(intervalId);
  }, [chatList]);

  useEffect(() => {
    if (isAutoScrollEnabled) {
      // Auto-scroll tiap kali ada perubahan pada chatList
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [displayResponse, chatList, isAutoScrollEnabled]);

  const handleScroll = () => {
    const container = chatContainerRef.current;
    if (!container) return;

    const isAtBottom =
      container.scrollHeight - container.scrollTop === container.clientHeight;

    setIsAutoScrollEnabled(isAtBottom);
  };

  useEffect(() => {
    if (chatList.some((chat) => chat.type === "error")) {
      setIsError(true);
    } else {
      setIsError(false);
    }
  }, [chatList]);

  return (
    <main
      ref={chatContainerRef}
      className="flex-1 overflow-y-auto"
      onScroll={handleScroll}
      style={{ scrollbarGutter: "stable both-edges" }}
    >
      <div className="py-4">
        <div className="flex justify-center">
          <div
            id="chat-list"
            className="chat w-full max-w-screen-md flex flex-col space-y-4"
          >
            {chatList.map((chat, chatId) =>
              chat.role === "user" ? (
                <UserChat key={chatId} chat={chat} />
              ) : (
                <AssistantChat
                  key={chatId}
                  chat={chat}
                  chatId={chatId}
                  chatListLength={chatList.length}
                  isCompletedTyping={isCompletedTyping}
                  displayResponse={displayResponse}
                  reloadChat={() => reloadChat()}
                />
              )
            )}
            <div ref={bottomRef} />
          </div>
        </div>
      </div>
    </main>
  );
};

export default ChatList;
