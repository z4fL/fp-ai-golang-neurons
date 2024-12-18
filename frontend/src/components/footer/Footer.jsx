import UploadSVG from "../svg/UploadSVG";
import SendSVG from "../svg/SendSVG";

const Footer = ({
  setIsModalOpen,
  query,
  setQuery,
  getResponse,
  isLoading,
  isError,
  errorType
}) => {
  return (
    <footer className="w-full max-w-screen-md p-4 mb-4 mx-auto bg-gray-200 dark:bg-gray-700 rounded-lg flex justify-center">
      <div className="w-full flex items-center space-x-2">
        {/* Upload File Button */}
        <button
          className="p-2 bg-lime-200 dark:bg-lime-300 rounded-md hover:bg-lime-300 dark:hover:bg-lime-400"
          title="Upload File"
          onClick={() => setIsModalOpen(true)}
          disabled={isLoading || isError}
        >
          <UploadSVG />
        </button>

        {/* Input Field */}
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Type here..."
          className="flex-1 p-2 dark:text-white dark:caret-white dark:bg-slate-600 border-0 rounded-md focus:outline-none focus:ring-2 focus:ring-lime-500 dark:focus:ring-lime-600 placeholder:text-slate-900 dark:placeholder:text-slate-300"
          onKeyDown={(e) => {
            if (e.key === "Enter" && query.trim() && !isLoading && !isError && errorType !== "text") {
              getResponse();
            }
          }}
        />

        {/* Send Button */}
        <button
          className={`bg-lime-400 px-3 py-2 rounded-md hover:bg-lime-500 ${
            isLoading && "disabled:bg-lime-700"
          }`}
          onClick={() => getResponse()}
          disabled={!query.trim() || isLoading || isError || errorType === "text"}
        >
          <SendSVG />
        </button>
      </div>
    </footer>
  );
};

export default Footer;
