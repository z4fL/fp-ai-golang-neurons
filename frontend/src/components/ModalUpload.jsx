import React, { useState, useEffect } from "react";

const ModalUpload = ({ isOpen, onClose, getResponse, file, setFile }) => {
  const [isFocused, setIsFocused] = useState(false); // state untuk kontrol fokus area drop
  const [isFileValid, setIsFileValid] = useState(true);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
  };

  const handleRemoveFile = () => {
    setFile(null);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    const draggedFile = e.dataTransfer.items[0];

    if (
      draggedFile &&
      draggedFile.kind === "file" &&
      draggedFile.type === "text/csv"
    ) {
      setIsFileValid(true);
    } else {
      setIsFileValid(false);
    }

    setIsFocused(true);
  };

  const handleDragLeave = () => {
    setIsFocused(false);
    setIsFileValid(true);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    const droppedFile = e.dataTransfer.files[0];

    if (droppedFile && droppedFile.type === "text/csv") {
      setFile(droppedFile);
    } else {
      console.log("only .csv");
    }

    setIsFileValid(true); // Reset validasi
    setIsFocused(false); // Hapus fokus setelah file di-drop
  };

  const handleCloseModal = () => {
    onClose();
  };

  const handleUploadFile = () => {
    getResponse();

    handleCloseModal();
  };

  return (
    isOpen && (
      <div className="fixed inset-0 bg-black bg-opacity-25 flex justify-center items-center z-50">
        <div className="bg-white dark:bg-gray-700 rounded-md shadow-xl w-full max-w-screen-md p-6 flex flex-col">
          {/* Drag and Drop Area */}
          {!file && (
            <div
              className={`border-dashed border-2 py-12 flex flex-col justify-center items-center
              ${
                isFocused
                  ? isFileValid
                    ? "dark:text-gray-800 border-lime-500 dark:border-lime-400 bg-lime-100 dark:bg-transparent"
                    : "border-red-500 bg-red-100 dark:bg-transparent"
                  : "border-gray-400 dark:border-gray-900"
              }`}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
            >
              <p className="text-gray-700 dark:text-gray-50">
                Drag and drop your file here or click the button below
              </p>
              <p className="mb-3 text-sm text-gray-600 dark:text-gray-200">
                Only .csv file can be uploaded
              </p>
              <input
                type="file"
                onChange={handleFileChange}
                className="hidden"
                id="file-input"
                name="file"
                accept="text/csv"
              />
              <label
                htmlFor="file-input"
                className="px-4 py-2 bg-lime-600 dark:bg-lime-500 text-white dark:text-gray-900 rounded-md cursor-pointer hover:bg-lime-500 dark:hover:bg-lime-600"
              >
                Upload File
              </label>
            </div>
          )}

          {/* Display selected file */}
          {file && (
            <div className="mt-4 bg-gray-100 dark:bg-gray-800 p-3 rounded-md flex justify-between items-center">
              <div className="flex flex-col dark:text-gray-100">
                <span className="text-ellipsis overflow-hidden">
                  {file.name}
                </span>
                <span>{(file.size / 1024).toFixed(2)} KB</span>
              </div>
              <button
                onClick={handleRemoveFile}
                className="text-red-500 hover:text-red-400"
              >
                Remove
              </button>
            </div>
          )}

          {/* Footer Buttons */}
          <div className="mt-6 flex justify-end space-x-4">
            <button
              onClick={handleUploadFile}
              className={`px-4 py-2 bg-lime-600 dark:bg-lime-400 dark:hover:bg-lime-500 text-white dark:text-gray-900 rounded-md hover:bg-lime-500 ${
                !file && "disabled:bg-lime-700 dark:disabled:bg-lime-600"
              }`}
              disabled={file ? false : true}
            >
              Upload and Analyze
            </button>
            <button
              onClick={() => {
                handleCloseModal();
                setFile(null);
              }}
              className="px-4 py-2 bg-gray-300 dark:bg-gray-400 text-gray-800 dark:text-gray-100 rounded-md hover:bg-gray-200 dark:hover:text-gray-700"
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
