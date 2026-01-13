import React, { useState, useEffect, useRef } from "react";
import { Card, CardContent, Box } from "@mui/material";
import TTS from "./tts";
import MarkDown from "./MarkDown";
import AudioControls from "./feedback";
import Image from "next/image";

const AutoSizingCard = ({ mdText = "", loading, error = null }) => {
  let requestId = null;
  if (typeof window !== 'undefined' && localStorage.getItem('sid')) {
    requestId = localStorage.getItem('sid')
  }

  const [fetched, setFetched] = useState(false);
  const [startTTS, setStartTTS] = useState(false);
  const [audioUrl, setAudioUrl] = useState(null);
  const [speechMarksData, setSpeechMarksData] = useState([]);
  const [markDownPrinted, setMarkDownPrinted] = useState(false);
  const [displayedText, setDisplayedText] = useState("");
  const [currentRequestId, setCurrentRequestId] = useState(null);
  const [hasError, setHasError] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const activePromiseRef = useRef(null);
  const isMountedRef = useRef(true);
  const timeoutIdRef = useRef(null);

  const waitForLocalStorage = (key, timeout = 30000) => {
    const startTime = Date.now();
    const promiseId = Symbol('promiseId');
    activePromiseRef.current = promiseId;

    return new Promise((resolve) => {
      // Use a local variable to track the timeout ID for this specific promise
      let currentTimeoutId = null;
      
      const checkItem = () => {
        // Check if this promise is still the active one
        if (activePromiseRef.current !== promiseId) {
          // Promise was cancelled, resolve with null
          if (currentTimeoutId) clearTimeout(currentTimeoutId);
          resolve(null);
          return;
        }

        // Check if component is still mounted
        if (!isMountedRef.current) {
          // Component unmounted, resolve with null
          if (currentTimeoutId) clearTimeout(currentTimeoutId);
          resolve(null);
          return;
        }

        const item = localStorage.getItem(key);
        if (item !== null) {
          if (activePromiseRef.current === promiseId && isMountedRef.current) {
            if (currentTimeoutId) clearTimeout(currentTimeoutId);
            resolve(item);
          } else {
            if (currentTimeoutId) clearTimeout(currentTimeoutId);
            resolve(null);
          }
        } else if (Date.now() - startTime >= timeout) {
          // Timeout reached - resolve with error object instead of rejecting
          if (currentTimeoutId) clearTimeout(currentTimeoutId);
          
          if (activePromiseRef.current === promiseId && isMountedRef.current) {
            // Still active, resolve with error object with user-friendly message
            resolve({ error: true, message: "Unable to process your request. The server may be temporarily unavailable." });
          } else {
            // No longer active, resolve with null
            resolve(null);
          }
        } else {
          currentTimeoutId = setTimeout(checkItem, 100);
          timeoutIdRef.current = currentTimeoutId;
        }
      };
      
      checkItem();
    });
  };

  useEffect(() => {
    isMountedRef.current = true;
    
    return () => {
      isMountedRef.current = false;
      activePromiseRef.current = null;
      if (timeoutIdRef.current) {
        clearTimeout(timeoutIdRef.current);
        timeoutIdRef.current = null;
      }
    };
  }, []);

  useEffect(() => {
    // Reset states when new request is made
    if (requestId && requestId !== currentRequestId) {
      // Cancel any previous promise
      activePromiseRef.current = null;
      if (timeoutIdRef.current) {
        clearTimeout(timeoutIdRef.current);
        timeoutIdRef.current = null;
      }
      
      setFetched(false);
      setStartTTS(false);
      setAudioUrl(null);
      setSpeechMarksData([]);
      setDisplayedText("");
      setCurrentRequestId(requestId);
      setHasError(false);
      setErrorMessage(null);
    }

    const fetchData = async () => {
      try {
        if (requestId) {
          const result = await waitForLocalStorage(requestId, 30000);
          
          // Check if component is still mounted and this is still the active request
          if (!isMountedRef.current || activePromiseRef.current === null) {
            return;
          }
          
          // If result is null, it means the promise was cancelled
          if (result === null) {
            return;
          }
          
          // Check if result is an error object
          if (result && typeof result === 'object' && result.error === true) {
            // This is a timeout error
            console.error("Error fetching data:", result.message);
            setHasError(true);
            setErrorMessage(result.message || "Your request couldn't be processed");
            setFetched(true);
            return;
          }
          
          // Result is a string (the data from localStorage)
          const data = JSON.parse(result);
          if (data) {
            setAudioUrl(data.audioData);
            setSpeechMarksData(data.speechMarksData);
          }
        }
        if (isMountedRef.current && activePromiseRef.current !== null) {
          setFetched(true);
        }
      } catch (error) {
        // Only handle error if component is still mounted and this is still the active request
        if (!isMountedRef.current || activePromiseRef.current === null) {
          return;
        }
        
        console.error("Error fetching data:", error);
        setHasError(true);
        setErrorMessage(error.message || "Your request couldn't be processed");
        setFetched(true);
      }
    };

    if (requestId) {
      fetchData();
    }
  }, [requestId]);

  // Check for error conditions
  useEffect(() => {
    if (error) {
      setHasError(true);
      setErrorMessage(error);
    } else if (!loading && !mdText && !hasError) {
      // If not loading, no text, and no explicit error yet, check if we should show error
      // This handles cases where backend didn't respond
      const timeoutId = setTimeout(() => {
        if (!mdText && !loading && isMountedRef.current) {
          setHasError(true);
          setErrorMessage("Your request couldn't be processed");
        }
      }, 5000); // Wait 5 seconds after loading stops to show error
      
      return () => clearTimeout(timeoutId);
    } else if (mdText && hasError) {
      // If we have text, clear error
      setHasError(false);
      setErrorMessage(null);
    }
  }, [error, loading, mdText, hasError]);

  const chunkLength = 100;
  const speed = 50;
  const [chunks, setChunks] = useState([]);

  const splitTextIntoChunks = (mdText, maxLength) => {
    const chunks = [];
    let start = 0;
    while (mdText && start < mdText.length) {
      chunks.push(mdText.slice(start, start + maxLength));
      start += maxLength;
    }
    return chunks;
  };

  useEffect(() => {
    setMarkDownPrinted(false);
    setChunks(splitTextIntoChunks(mdText, chunkLength));
  }, [mdText]);

  useEffect(() => {
    if (chunks.length === 0) return;
    let index = -1;
    let timerId;
    const typeNextChunk = () => {
      if (index < chunks.length - 1) {
        setDisplayedText((prev) => prev + chunks[index]);
        index++;
        timerId = setTimeout(typeNextChunk, speed);
      }
      if (index === chunks.length - 1) {
        setMarkDownPrinted(true);
      }
    };
    typeNextChunk();
    return () => clearTimeout(timerId);
  }, [chunks]);

  // Error display component
  const ErrorDisplay = () => (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        padding: "60px 20px",
        minHeight: "500px",
        maxWidth: "900px",
        margin: "20px auto",
      }}
    >
      <div
        style={{
          width: "200px",
          height: "200px",
          marginBottom: "20px",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          // backgroundColor: "#2a2a2a",
          borderRadius: "16px",
        }}
      >
        {/* You can replace this with an actual error image/GIF */}
        <div
          style={{
            fontSize: "80px",
            color: "#f1c40f",
          }}
        >
          ⚠️
        </div>
        {/* Or use an image:
        <Image
          src="/error.gif" // Add your error image/GIF to public folder
          alt="Error"
          width={300}
          height={300}
          style={{ objectFit: "contain" }}
        />
        */}
      </div>
      <p
        style={{
          color: "#f2bb4b",
          fontSize: "18px",
          textAlign: "center",
          marginTop: "20px",
        }}
      >
        {errorMessage || "Your request couldn't be processed"}
      </p>
      <p
        style={{
          color: "#888",
          fontSize: "14px",
          textAlign: "center",
          marginTop: "10px",
        }}
      >
        Ultra-rare-sophisticated occurrence. Please try again.
      </p>
    </div>
  );

  // If there's an error or no content after loading, show error display directly (not in Card)
  if (hasError || (!mdText && !loading)) {
    return (
      <div
        style={{
          padding: "1px",
          position: "relative",
          maxWidth: 2000,
          margin: "20px auto",
          borderRadius: "24px",
        }}
      >
        <div
          style={{
            position: "absolute",
            top: "-1px",
            left: "-1px",
            right: "-1px",
            bottom: "-1px",
            zIndex: -1,
            // background: "linear-gradient(45deg, #3498db, #f1c40f)",
            filter: "blur(15px)",
            borderRadius: "30px",
          }}
        />
        <ErrorDisplay />
      </div>
    );
  }

  return (
    <div
      style={{
        padding: "1px",
        position: "relative",
        maxWidth: 2000,
        margin: "20px auto",
        borderRadius: "24px",
      }}
    >
      <div
        style={{
          position: "absolute",
          top: "-1px",
          left: "-1px",
          right: "-1px",
          bottom: "-1px",
          zIndex: -1,
          background: "linear-gradient(45deg, #3498db, #f1c40f)",
          filter: "blur(15px)",
          borderRadius: "30px",
        }}
      />

      <Card
        className="max-w-4xl mx-auto p-6"
        style={{
          borderRadius: "24px",
          background: "#1e1e1e",
          color: "#f2bb4b",
        }}
      >
        {loading ? (
          <CardContent
            sx={{
              width: { xs: "250px", sm: "2000px" },
              height: "500px",
            }}
          >
            <Box
              sx={{
                width: 15,
                height: 15,
                borderRadius: "100%",
                backgroundColor: "#f1c40f",
                animation: "pulse 1s infinite",
                "@keyframes pulse": {
                  "0%": { transform: "scale(1)" },
                  "50%": { transform: "scale(1.3)" },
                  "100%": { transform: "scale(1)" },
                },
              }}
            />
          </CardContent>
        ) : (
          <CardContent
            sx={{
              paddingTop: 0,
              maxWidth: { xs: "450px", sm: "900px" },
              maxHeight: { xs: "450px", sm: "600px" },
              overflowY: "auto",
            }}
          >
            {fetched === false ? (
              <MarkDown text={displayedText} />
            ) : (
              <>
                <div className="sticky top-0 bg-[#1e1e1e] z-10">
                  <AudioControls startTTS={startTTS} setStartTTS={setStartTTS} />
                </div>
                <TTS
                  text={mdText}
                  audioUrl={audioUrl}
                  speechMarksData={speechMarksData}
                  startTTS={startTTS}
                  setStartTTS={setStartTTS}
                />
              </>
            )}
          </CardContent>
        )}
      </Card>
    </div>
  );
};

export default AutoSizingCard;