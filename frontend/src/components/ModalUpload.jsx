import React, { useState } from "react";

const ModalUpload = ({ isOpen, onClose }) => {
  const [file, setFile] = useState(null);
  const [isFocused, setIsFocused] = useState(false); // state untuk kontrol fokus area drop

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
  };

  const handleRemoveFile = () => {
    setFile(null);
  };

  const handleDragOver = (e) => {
    e.preventDefault();

    setIsFocused(true);
  };

  const handleDragLeave = () => {
    setIsFocused(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    const droppedFile = e.dataTransfer.files[0];
    setFile(droppedFile);
    setIsFocused(false); // Hapus fokus setelah file di-drop
  };

  const handleCloseModal = () => {
    onClose();
    handleRemoveFile();
  };

  return (
    isOpen && (
      <div className="fixed inset-0 bg-black bg-opacity-25 flex justify-center items-center z-50">
        <div className="bg-white rounded-md shadow-xl w-full max-w-screen-md p-6 flex flex-col">
          {/* Drag and Drop Area */}
          <div
            className={`border-dashed border-2 py-12 flex flex-col justify-center items-center
              ${isFocused ? "border-lime-500 bg-lime-100" : "border-gray-400"}`}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
          >
            <p className="mb-3 text-gray-700">
              Drag and drop your file here or click the button below
            </p>
            <input
              type="file"
              onChange={handleFileChange}
              className="hidden"
              id="file-input"
              name="file"
            />
            <label
              htmlFor="file-input"
              className="px-4 py-2 bg-lime-600 text-white rounded-md cursor-pointer hover:bg-lime-500"
            >
              Upload File
            </label>
          </div>

          {/* Display selected file */}
          {file && (
            <div className="mt-4 bg-gray-100 p-3 rounded-md flex justify-between items-center">
              <span>{file.name}</span>
              <button
                onClick={handleRemoveFile}
                className="text-red-500 hover:underline"
              >
                Remove
              </button>
            </div>
          )}

          {/* Footer Buttons */}
          <div className="mt-6 flex justify-end space-x-4">
            <button
              onClick={() => console.log("File uploaded: ", file)}
              className="px-4 py-2 bg-lime-600 text-white rounded-md hover:bg-lime-500"
            >
              Upload Now
            </button>
            <button
              onClick={handleCloseModal}
              className="px-4 py-2 bg-gray-300 text-gray-800 rounded-md hover:bg-gray-400"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    )
  );
};

export default ModalUpload;
