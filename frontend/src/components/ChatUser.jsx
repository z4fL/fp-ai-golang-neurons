import React from "react";

const ChatUser = ({ type, content }) => {
  return (
    <div>
      {type === "text" ? (
        <p>{content}</p>
      ) : (
        <>
          <p>{content.name}</p>
          <p className="text-sm text-slate-200">
            {(content.size / 1024).toFixed(2)} KB
          </p>
        </>
      )}
    </div>
  );
};

export default ChatUser;
