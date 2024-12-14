import { useState } from "react";
import Markdown from "react-markdown";
import AnimateSpinSVG from "./components/svg/AnimateSpinSVG";
import ErrorChat from "./components/Chat/ErrorChat";
import ReloadSVG from "./components/svg/ReloadSVG";

const ChatArea = () => {
  const [chatHistory, setChatHistory] = useState([
    {
      content: "Hello, how can I help you?",
      id: 1,
      role: "assistant",
    },
    {
      content: {
        name: "data-series.csv",
        size: 1404,
      },
      id: 2,
      role: "user",
      type: "file",
    },
    {
      id: 3,
      role: "assistant",
      content:
        "From the provided data, here are the Least Electricity: TV and the Most Electricity: EVCar.",
      type: "text",
    },
    {
      content: "Can unplugging the TV completely save energy?",
      id: 4,
      role: "user",
      type: "text",
    },
    {
      content:
        "Yes, unplugging the TV when not in use can save energy. Even when turned off, televisions can continue to draw power, a phenomenon known as \"phantom load\" or \"vampire power.\" By unplugging your TV, you eliminate this unnecessary energy consumption.\n\nHowever, it's worth noting that the amount of energy saved depends on the TV's power rating and how often it's used. For instance, if your TV consumes 100 watts and you watch it for 5 hours a day, that's 500 watt-hours or 1.5 kWh per day. Over a month, that adds up to about 45 kWh. If your electricity rate is $0.10 per kWh, you're spending about $4.50 per month on TV power. By unplugging the TV, you could save that amount.\n\nMoreover, unplugging your TV can also reduce your carbon footprint, as less energy consumption means fewer greenhouse gas emissions, assuming your electricity comes from fossil fuels.\n\nRemember, though, to only unplug your TV when it's not in use. If you're planning to use it again soon, plug it back in. Also, consider using a power strip with an on/off switch for multiple devices, which can make it easier to cut power to several things at once.\n\nLastly, for TVs that have remote controls or built-in digital tuners, consider using power-saving features if available. Some modern TVs have an \"auto-off\" feature that turns the TV off after a certain period of inactivity.",
      id: 5,
      role: "assistant",
      type: "text",
    },
  ]);

  return (
    <div
      id="chat-list"
      className="chat w-full max-w-screen-md flex flex-col space-y-4"
    >
      {chatHistory.map((chat, chatId) =>
        chat.role === "user" ? (
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
        ) : (
          <div
            className={`self-start p-3 rounded-md bg-slate-100 text-slate-900`}
          >
            <div className="prose prose-base">
              <Markdown>{chat.content}</Markdown>
            </div>
          </div>
        )
      )}
    </div>
  );
};

export default ChatArea;
