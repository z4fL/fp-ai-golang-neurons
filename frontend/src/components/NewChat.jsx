import Markdown from "react-markdown";

const NewChat = () => {
  return (
    <main
      className="flex-1 overflow-y-auto"
      style={{ scrollbarGutter: "stable both-edges" }}
    >
      <div className="py-4">
        <div className="flex justify-center">
          <div
            id="chat-list"
            className="chat w-full max-w-screen-md flex flex-col space-y-4"
          >
            <div className="self-start p-3 rounded-md bg-slate-100 text-slate-900">
              <div className="prose prose-base">
                <Markdown>Hello, how can I help you?</Markdown>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
};

export default NewChat;
