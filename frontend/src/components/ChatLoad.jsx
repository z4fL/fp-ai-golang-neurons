
const ChatLoad = () => {
  return (
    <main className="flex-1 flex justify-center items-center">
      <div
        id="chat-list"
        className="chat w-full max-w-screen-md flex flex-col space-y-4 items-center justify-center"
      >
        <img
          src="/logo.png"
          alt="Logo"
          className="animate-bounce max-w-20 md:max-w-full h-auto"
        />
      </div>
    </main>
  );
};

export default ChatLoad;
