import React, { useState, useEffect } from 'react';
import FeedbackDrawer from './FeedbackDrawer';

// Feedback Component
const Feedback = ({ open, feedback, setFeedback }) => {
    const [isAnimating, setIsAnimating] = useState(false);
    const handleThumbsUp = () => {
        setIsAnimating(true);
        setTimeout(() => setIsAnimating(false), 500);
        const requestId = localStorage.getItem('sid')
        const existingFeedback = localStorage.getItem(`feedback_${requestId}`);
        if (existingFeedback) {
            const feedbackObj = JSON.parse(existingFeedback);
            const updatedFeedback = {
                ...feedbackObj,
                feedback: 1,
            };
            localStorage.setItem(`feedback_${requestId}`, JSON.stringify(updatedFeedback));
        }
        setFeedback(1)
    }
    const handleThumbsDown = () => {
        open(true);
    }

    return (
        <div className="flex gap-1">
            <button
                onClick={handleThumbsUp}
                className="p-1 text-gray-500 hover:text-gray-700 transition-colors duration-200 rounded-full hover:bg-[#454545]"
            >
                <svg
                    className={`${isAnimating ? "animate-pop" : ""}`}
                    viewBox="0 0 24 24"
                    fill={feedback === 1 ? "#f1c40f" : "#1e1e1e"}
                    stroke="#f2bb4b"
                    strokeWidth="2"
                    width="16"
                    height="16"
                >
                    <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3" />
                </svg>
            </button>
            <button
                onClick={handleThumbsDown}
                className="p-1 text-gray-500 hover:text-gray-700 transition-colors duration-200 rounded-full hover:bg-[#454545]"
            >
                <svg width="16" height="16" viewBox="0 0 24 24" fill={feedback === -1 ? "#f1c40f" : "#1e1e1e"} stroke="#f2bb4b" strokeWidth="2">
                    <path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zm7-13h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3" />
                </svg>
            </button>
        </div>
    );
}

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
        <div className="relative h-8">
            <div
                className={`absolute inset-0 w-full h-full rounded-full 
        ${isRendered ? 'border-2 border-[#f2bb4b] animate-ping' : ''}`}
                style={{ pointerEvents: 'none' }}
            />
            <button
                onClick={handlePlayPause}
                className="button cursor-pointer flex rounded-full leading-5 items-center text-xs px-3 py-1 border border-[#333333] bg-transparent hover:bg-[#454545] transition-colors duration-200"
            >
                {startTTS ? (
                    <>
                        <svg className="mr-1.5" width="8" height="8" viewBox="0 0 8 8" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <rect width="8" height="8" rx="1.5" fill="currentColor" />
                        </svg>
                        Stop
                    </>
                ) : (
                    <>
                        <svg className="mr-1.5" width="8" height="10" viewBox="0 0 8 10" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M1 9V1L7 5.4L1 9Z" fill="currentColor" stroke="currentColor" />
                        </svg>
                        Listen
                    </>
                )}
            </button>
        </div>
    );
};

// Parent Component
const AudioControls = ({ startTTS, setStartTTS }) => {
    const [isOpen, setIsOpen] = useState(false);
    const [feedback, setFeedback] = useState(0);

    const handleFeedbackSubmit = () => {
        setFeedback(-1);
    }
    return (
        <div>
            <FeedbackDrawer
                open={isOpen}
                onClose={() => setIsOpen(false)}
                onSubmit={handleFeedbackSubmit}
            />
            <div className="w-full h-12 flex justify-between items-center">
                <div className="z-10">
                    <PlayButton startTTS={startTTS} setStartTTS={setStartTTS} />
                </div>
                <div>
                    <Feedback open={setIsOpen} feedback={feedback} setFeedback={setFeedback} />
                </div>
                {/* Horizontal Line */}
            </div>
            <hr className="w-full h-[1px] bg-[#333333] border-none" />

        </div>
    );
};

export default AudioControls;