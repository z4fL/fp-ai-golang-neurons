import React from "react";
import Markdown from "react-markdown";
import AnimateSpinSVG from "../svg/AnimateSpinSVG";
import CursorSVG from "../svg/CursorSVG";
import ErrorChat from "./ErrorChat";
import ReloadSVG from "../svg/ReloadSVG";

const AssistantChat = ({
  chat,
  chatListLength,
  isCompletedTyping,
  displayResponse,
  reloadChat,
}) => {
  const baseClass = "self-start p-3 rounded-md";
  const contentClass =
    chat.type === "error"
      ? "bg-red-200 border-2 border-red-400 text-slate-900"
      : "bg-slate-100 text-slate-900";

  if (chat.id !== chatListLength - 1) {
    return (
      <div className={`${baseClass} ${contentClass}`}>
        <div className="prose prose-base">
          <Markdown>{chat.content}</Markdown>
        </div>
      </div>
    );
  }

  if (chat.content === "LOADING...") {
    return (
      <div className={`${baseClass} ${contentClass}`}>
        <div className="prose prose-base flex items-center">
          <AnimateSpinSVG className="-ml-1 mr-3 h-5 w-5 text-slate-950" />
          {chat.content}
        </div>
      </div>
    );
  }

  if (chat.type === "error") {
    return (
      <div className={`flex items-center p-3`}>
        <ErrorChat
          baseClass={`${baseClass} ${contentClass}`}
          content={chat.content}
        />
        <div
          className="ml-2 p-2 transition-transform duration-300 hover:rotate-[270deg]"
          onClick={() => reloadChat()}
        >
          <ReloadSVG />
        </div>
      </div>
    );
  }

  return (
    <div className={`${baseClass} ${contentClass}`}>
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
    </div>
  );
};

export default AssistantChat;
