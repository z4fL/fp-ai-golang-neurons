import React from "react";

const ErrorChat = ({baseClass, content}) => {
  return (
    <div className={`${baseClass}`}>
      <div className="prose prose-base flex items-center text-slate-900 dark:text-slate-200">{content}</div>
    </div>
  );
};

export default ErrorChat;
