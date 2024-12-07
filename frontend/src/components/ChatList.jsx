import React, { useState, useRef } from "react";
import { useEffect } from "react";
import CursorSVG from "./svg/Cursor";
import Markdown from "react-markdown";

const ChatList = ({ chatList }) => {
  const [displayResponse, setDisplayResponse] = useState("");
  const [isCompletedTyping, setIsCompletedTyping] = useState(false);

  const [isAutoScrollEnabled, setIsAutoScrollEnabled] = useState(true);
  const bottomRef = useRef(null);
  const chatContainerRef = useRef(null);

  useEffect(() => {
    if (!chatList.length) {
      return;
    }

    setIsCompletedTyping(false);

    let i = 0;
    const responseAssistant = chatList[chatList.length - 1].content;

    const intervalId = setInterval(() => {
      setDisplayResponse(responseAssistant.slice(0, i));
      i++;

      if (i > responseAssistant.length) {
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

    setIsAutoScrollEnabled(isAtBottom)
  };

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
            {chatList.map((chat, chatId) => (
              <div
                key={chatId}
                className={`p-3 rounded-md ${
                  chat.role === "assistant"
                    ? "self-start bg-slate-100 text-slate-900"
                    : "self-end bg-gray-800 text-white"
                }`}
              >
                {chat.role === "user" && (
                  <div>
                    {chat.type === "text" ? (
                      <p>{chat.content}</p>
                    ) : (
                      <>
                        <p>{chat.content.name}</p>
                        <p className="text-sm text-slate-200">
                          {(chat.content.size / 1024).toFixed(2)} KB
                        </p>
                      </>
                    )}
                  </div>
                )}

                {chat.role === "assistant" &&
                  chatId !== chatList.length - 1 && (
                    <div className="prose prose-base">
                      <Markdown>{chat.content}</Markdown>
                    </div>
                  )}

                {chat.role === "assistant" &&
                  chatId === chatList.length - 1 && (
                    <div className="prose prose-base">
                      {!isCompletedTyping ? (
                        <>
                          {displayResponse}
                          <CursorSVG />
                        </>
                      ) : (
                        <Markdown>{displayResponse}</Markdown>
                      )}
                    </div>
                  )}
              </div>
            ))}
            <div ref={bottomRef} />
          </div>
        </div>
      </div>
    </main>
  );
};

export default ChatList;
