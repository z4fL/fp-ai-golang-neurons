import React from "react";

const ChatUser = ({ type, value }) => {
  return (
    <div>
      {type === "text" ? (
        <p>{value.content}</p>
      ) : (
        <>
          <p>{value.name}</p>
          <p className="text-sm text-slate-200">
            {(value.size / 1024).toFixed(2)} KB
          </p>
        </>
      )}
    </div>
  );
};

export default ChatUser;
