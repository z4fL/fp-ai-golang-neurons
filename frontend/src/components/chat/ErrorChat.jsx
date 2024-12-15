import React from "react";

const ErrorChat = ({baseClass, content}) => {
  return (
    <div className={`${baseClass}`}>
      <div className="prose prose-base">{content}</div>
    </div>
  );
};

export default ErrorChat;
