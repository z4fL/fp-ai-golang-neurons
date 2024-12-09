import React from "react";
import UploadSVG from "./svg/UploadSVG";
import SendSVG from "./svg/SendSVG";

const Footer = ({
  setIsModalOpen,
  query,
  setQuery,
  getResponse,
  isLoading,
  isError,
}) => {
  
  return (
    <footer className="w-full max-w-screen-md p-4 mb-4 mx-auto bg-gray-200 rounded-lg flex justify-center">
      <div className="w-full flex items-center space-x-2">
        {/* Upload File Button */}
        <button
          className="p-2 bg-lime-200 rounded-md hover:bg-lime-300"
          title="Upload File"
          onClick={() => setIsModalOpen(true)}
          disabled={isLoading}
        >
          <UploadSVG />
        </button>

        {/* Input Field */}
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Ketik pesan..."
          className="flex-1 p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-lime-500 placeholder:text-slate-500"
          onKeyDown={(e) => {
            if (e.key === "Enter" && query.trim() && !isLoading && !isError) {
              getResponse("text");
            }
          }}
        />

        {/* Send Button */}
        <button
          className={`bg-lime-400 px-3 py-2 rounded-md hover:bg-lime-500 ${
            isLoading && "disabled:bg-lime-700"
          }`}
          onClick={() => getResponse("text")}
          disabled={!query.trim() || isLoading || isError}
        >
          <SendSVG />
        </button>
      </div>
    </footer>
  );
};

export default Footer;
