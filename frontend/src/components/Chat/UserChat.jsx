import React from "react";

const UserChat = ({ chat }) => (
  <div className="self-end bg-gray-800 text-white p-3 rounded-md">
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
);

export default UserChat;
