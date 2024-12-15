import React from "react";

const UserChat = ({ chat }) => (
  <div className="self-end bg-gray-300 dark:bg-gray-950 text-gray-800 dark:text-gray-200 p-4 rounded-md">
    {chat.type === "text" ? (
      <p>{chat.content}</p>
    ) : (
      <>
        <p>{chat.content.name}</p>
        <p className="text-sm text-gray-800 dark:text-gray-500">
          {(chat.content.size / 1024).toFixed(2)} KB
        </p>
      </>
    )}
  </div>
);

export default UserChat;
