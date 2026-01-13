import { v4 as uuidv4 } from 'uuid';
import { getCachedData, setCachedData } from './cacheUtils';
import { formatter } from './timeFormatter';

const wsbackendUrl = process.env.NEXT_PUBLIC_WS_URL;

export const processContent = async (url, onSummarize, onTTS, onError) => {
  try {
    const summaryCacheKey = `summary_${url}`;

    const socket = new WebSocket(wsbackendUrl);
    const requestId = uuidv4();
    localStorage.setItem("sid", requestId);
    localStorage.setItem(`feedback_${requestId}`, JSON.stringify({
      feedback: 0,
      feedbackText: "",
      ttsplayedDuration: 0,
      navigatedToClean: false,
    }));

    socket.onopen = () => {
      socket.send(
        JSON.stringify({
          type: "process",
          requestId: requestId,
          payload: { url },
        })
      );
    };

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.requestId !== requestId) {
          console.warn("[WebSocket] Received unrelated response, ignoring.");
          return;
        }
        if (message.status === "error") {
          console.error("[WebSocket] Error message received:", message.message);
          if (onError) onError(message.message);
          return;
        }
        if (message.type === "summarize" && onSummarize) {
          setCachedData(summaryCacheKey, message.data);
          onSummarize(message.data);
        }
        if (message.type === "tts") {
          localStorage.setItem(message.requestId, JSON.stringify(message.data));
        }
      } catch (error) {
        console.error("[WebSocket] Failed to parse WebSocket message", error);
        if (onError) onError("Failed to parse WebSocket message");
      }
    };

    socket.onerror = (error) => {
      console.error("[WebSocket] WebSocket error:", error);
      if (onError) onError("WebSocket error");
    };

    socket.onclose = () => {
    };

    return () => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  } catch (error) {
    console.error("[WebSocket] Failed to establish WebSocket connection", error);
    if (onError) onError("Failed to establish WebSocket connection");
  }
};
