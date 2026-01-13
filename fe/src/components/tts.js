import Markdown from 'markdown-to-jsx';
import React, { useEffect, useRef, useState } from 'react';
import { baseMarkdownOptions } from '../utils/markdownOptions';

const ZWSP = '\u200B'; // Zero-Width Space

const TextToSpeech = ({ text, audioUrl, speechMarksData, startTTS, setStartTTS }) => {
  const audioRef = useRef(null);
  const [highlights, setHighlights] = useState([]);
  const [displayText, setDisplayText] = useState(text);
  const highlightTimeoutRef = useRef(null);
  const lastPositionRef = useRef(0);
  const hasEndedRef = useRef(false);
  const playPromiseRef = useRef(null);
  const maxPlayedDurationRef = useRef(0);

  useEffect(() => {
    const audio = audioRef.current;
  
    if (audio) {
      const handleTimeUpdate = () => {
        maxPlayedDurationRef.current = Math.max(maxPlayedDurationRef.current, audio.currentTime);
        const requestId = localStorage.getItem('sid')
        const existingFeedback = localStorage.getItem(`feedback_${requestId}`);
        if (existingFeedback) {
            // Parse existing feedback and update with new flag
            const feedbackObj = JSON.parse(existingFeedback);
            const updatedFeedback = {
                ...feedbackObj,
                ttsplayedDuration: maxPlayedDurationRef.current,
            };
            // Save updated feedback back to localStorage
            localStorage.setItem(`feedback_${requestId}`, JSON.stringify(updatedFeedback));
        }
      };
      audio.addEventListener("timeupdate", handleTimeUpdate);
      return () => {
        audio.removeEventListener("timeupdate", handleTimeUpdate);
      };
    }
  }, []);
  
  const markdownOptions = {
    ...baseMarkdownOptions,
    overrides: {
      ...baseMarkdownOptions.overrides,
      highlight: {
        component: ({ children }) => (
          <mark className="bg-amber-300 text-black px-[2px] rounded">{children}</mark>
        ),
      },
    },
  };

  const processTextHighlights = () => {
    const wordPositions = [];
    let textLength = text.length;
    let totalWords = speechMarksData.length;
    let i = 0;

    for (let j = 0; j < totalWords && i < textLength; j++) {
      let word = speechMarksData[j].value;
      if(j > 0 && word.includes(speechMarksData[j-1].value) && word.includes(" ")){
        word = word.split(" ")[1];
      }
      let match = word[0];
      while (i < textLength && text[i] !== match) {
        i++;
      }

      if (i < textLength) {
        wordPositions.push({
          start: i,
          end: i + word.length,
          time: speechMarksData[j].time,
          word: word
        });
      }
      i += word.length;
    }

    return wordPositions;
  };

  useEffect(() => {
    const positions = processTextHighlights();
    setHighlights(positions);
    setDisplayText(text);
  }, [text, speechMarksData]);

  useEffect(() => {
    if (!audioRef.current) return;

    const handlePlayPause = async () => {
      try {
        if (startTTS) {
          if (hasEndedRef.current) {
            lastPositionRef.current = 0;
            hasEndedRef.current = false;
          }
          // Handle play with promise
          if (playPromiseRef.current !== null) {
            await playPromiseRef.current;
          }
          playPromiseRef.current = audioRef.current.play();
          await playPromiseRef.current;
          startHighlighting();
        } else {
          // Handle pause
          if (playPromiseRef.current !== null) {
            await playPromiseRef.current;
          }
          audioRef.current.pause();
          lastPositionRef.current = audioRef.current.currentTime * 1000;
          clearHighlightTimeouts();
        }
      } catch (error) {
        if (error.name !== 'AbortError') {
          console.error('Audio playback error:', error);
        }
      } finally {
        playPromiseRef.current = null;
      }
    };

    handlePlayPause();

    return () => {
      clearHighlightTimeouts();
      if (audioRef.current) {
        audioRef.current.pause();
      }
    };
  }, [startTTS]);

  const clearHighlightTimeouts = () => {
    if (highlightTimeoutRef.current) {
      highlightTimeoutRef.current.forEach(timeout => clearTimeout(timeout));
      highlightTimeoutRef.current = [];
    }
  };

  const startHighlighting = () => {
    clearHighlightTimeouts();
    highlightTimeoutRef.current = [];
    const currentTime = lastPositionRef.current;
    let startIndex = highlights.findIndex(highlight => highlight.time >= currentTime);
    if (startIndex === -1) startIndex = 0;

    audioRef.current.currentTime = currentTime / 1000;

    highlights.slice(startIndex).forEach((highlight) => {
      const timeout = setTimeout(() => {
        const before = text.substring(0, highlight.start);
        const highlighted = text.substring(highlight.start, highlight.end);
        const after = text.substring(highlight.end);

        let prefix = '';
        // If the highlighted word is at the start of the text or follows a newline,
        // it's at the beginning of a block. Prepending a zero-width space
        // ensures that `markdown-to-jsx` treats the line as a paragraph
        // and keeps the highlighted element inside it, preventing layout shifts.
        if (highlight.start === 0 || text[highlight.start - 1] === '\n') {
          prefix = ZWSP;
        }

        const newText = `${before}${prefix}<highlight>${highlighted}</highlight>${after}`;
        setDisplayText(newText);
      }, highlight.time - currentTime);

      highlightTimeoutRef.current.push(timeout);
    });
  };

  const handleAudioEnd = () => {
    clearHighlightTimeouts();
    setDisplayText(text);
    hasEndedRef.current = true;
    lastPositionRef.current = 0;
    setStartTTS(false);
    playPromiseRef.current = null;
  };

  const handleTimeUpdate = () => {
    lastPositionRef.current = audioRef.current.currentTime * 1000;
  };

  return (
    <div>
      <div
        id="tts-text"
        className="leading-relaxed whitespace-pre-wrap"
      >
        <Markdown options={markdownOptions}>{displayText}</Markdown>
      </div>
      <audio
        ref={audioRef}
        src={audioUrl}
        onEnded={handleAudioEnd}
        onTimeUpdate={handleTimeUpdate}
        className="hidden"
      />
    </div>
  );
};

export default TextToSpeech;