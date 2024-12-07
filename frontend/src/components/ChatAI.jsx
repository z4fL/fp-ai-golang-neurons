import React from "react";
import Markdown from "react-markdown";

const ChatAI = ({ content }) => {

  return (
    <div className="prose prose-base">
      <Markdown>{content}</Markdown>
    </div>
  );
};

export default ChatAI;
