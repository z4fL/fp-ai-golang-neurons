import { useEffect, useRef, useState } from "react";
import { useLocation, useOutletContext } from "react-router";
import UserChat from "./chat/UserChat";
import AssistantChat from "./chat/AssistantChat";

const ChatArea = () => {
  const { chatHistory, setIsLoading, setIsError, reloadChat } =
    useOutletContext();

  const location = useLocation();

  const [displayResponse, setDisplayResponse] = useState("");
  const [isCompletedTyping, setIsCompletedTyping] = useState(false);
  const [isAutoScrollEnabled, setIsAutoScrollEnabled] = useState(true);

  const bottomRef = useRef(null);
  const chatContainerRef = useRef(null);

  useEffect(() => {
    const fromNavigate = location.state?.fromNavigate;
    window.history.replaceState({}, "");
    if (chatHistory.some((chat) => chat.type === "error")) {
      setIsError(true);
    } else {
      setIsError(false);
    }

    const lastChat = chatHistory[chatHistory.length - 1];

    if (lastChat.role !== "assistant") return;
    if (!fromNavigate) {
      setDisplayResponse(lastChat.content);
      setIsCompletedTyping(true);
      return;
    }

    console.log(fromNavigate);

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
  }, [chatHistory, location]);

  useEffect(() => {
    if (isAutoScrollEnabled) {
      // Auto-scroll tiap kali ada perubahan pada chatHistory
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [displayResponse, chatHistory, isAutoScrollEnabled]);

  const handleScroll = () => {
    const container = chatContainerRef.current;
    if (!container) return;

    const isAtBottom =
      container.scrollHeight - container.scrollTop === container.clientHeight;

    setIsAutoScrollEnabled(isAtBottom);
  };

  return (
    <main
      ref={chatContainerRef}
      className="flex-1 overflow-y-auto scrollbar"
      onScroll={handleScroll}
      style={{ scrollbarGutter: "stable both-edges" }}
    >
      <div className="py-4">
        <div className="flex justify-center">
          <div
            id="chat-list"
            className="chat w-full max-w-screen-md flex flex-col space-y-4"
          >
            {chatHistory.map((chat) =>
              chat.role === "user" ? (
                <UserChat key={chat.id} chat={chat} />
              ) : (
                <AssistantChat
                  key={chat.id}
                  chat={chat}
                  chatListLength={chatHistory.length}
                  isCompletedTyping={isCompletedTyping}
                  displayResponse={displayResponse}
                  reloadChat={reloadChat}
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

export default ChatArea;
