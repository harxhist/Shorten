import React, { useState, useEffect } from 'react';
const PlayButton = ({ startTTS, setStartTTS }) => {
  const [isRendered, setIsRendered] = useState(false);
  useEffect(() => {
    setIsRendered(true);
    const timer = setTimeout(() => setIsRendered(false), 1000);
    return () => clearTimeout(timer);
  }, []);
  const handlePlayPause = () => {
    setStartTTS(!startTTS);
  };
  return (
    <div className="relative w-10 h-10 m-5">
      <div
        className={`absolute inset-0 w-full h-full rounded-full 
          ${isRendered ? 'border-2 border-[#f2bb4b] animate-ping' : ''}`}
      />
      <div>
        {startTTS ? (<>
          <button onClick={handlePlayPause} class="button cursor-pointer flex rounded-full leading-5 flex-row items-center text-xs px-[11px] py-[2px] border  bg-transparent hover:bg-[#123e23]" type="button">
            <svg class="mr-1.5" width="8" height="8" viewBox="0 0 8 8" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="8" height="8" rx="1.5" fill="currentColor"></rect>
            </svg>
            Stop
          </button>
        </>) : (
          <>
            <button onClick={handlePlayPause} class="button cursor-pointer leading-5 rounded-full flex flex-row items-center text-xs px-[11px] py-[2px] border bg-transparent hover:bg-[#123e23]" type="button">
              <svg class="mr-1.5" width="8" height="10" viewBox="0 0 8 10" fill="none" xmlns="http://www.w3.org/2000/svg" >
                <path d="M1 9V1L7 5.4L1 9Z" fill="currentColor" stroke="currentColor"></path>
              </svg>
              Listen
            </button>
          </>
        )}
      </div>
    </div>
  );
};

export default PlayButton;