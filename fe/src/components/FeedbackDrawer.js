import React, { useState } from 'react';
import { Drawer, Button } from '@mui/material';
// import { Input } from '@mui/material';

const FeedbackDrawer = ({ open, onClose, onSubmit }) => {
    const [feedbacktxt, setFeedback] = useState('');

    const handleSubmit = () => {
        const requestId = localStorage.getItem('sid')
        const existingFeedback = localStorage.getItem(`feedback_${requestId}`);
        if (existingFeedback) {
            // Parse existing feedback and update with new flag
            const feedbackObj = JSON.parse(existingFeedback);
            const updatedFeedback = {
                ...feedbackObj,
                feedbackText: feedbacktxt,
                feedback: -1,
            };

            // Save updated feedback back to localStorage
            localStorage.setItem(`feedback_${requestId}`, JSON.stringify(updatedFeedback));
        }
        if (onSubmit) {
            onSubmit()
        }
        setFeedback('');
        onClose();
    };

    return (
        <Drawer
            PaperProps={{
                style: {
                    backgroundColor: '#1e1e1e'
                }
            }}
            anchor="bottom"
            open={open}
            onClose={onClose}
        >
            <div className="p-4 bg-[#1e1e1e] flex flex-row items-center gap-4 min-h-[80px] max-w-screen-lg mx-auto w-full">
                <input
                    value={feedbacktxt}
                    onChange={(e) => setFeedback(e.target.value)}
                    placeholder="(Optional) What could be better?"
                    maxLength={100}
                    autoFocus={false}
                    className="
        w-full
        flex-1
        px-4
        py-2
        text-[#f2bb4b]
        bg-neutral-900
        rounded-full
        placeholder-[#5c6370]
        focus:outline-none
        duration-200
        sm:text-sm
        md:text-base
        lg:px-6
        lg:py-2
        border-[#333333]
        border
      "
                />
                <button
                    onClick={handleSubmit}
                    className="button text-[#f2bb4b] cursor-pointer flex rounded-full leading-6 items-center text-base px-6 py-2 border border-[#333333] bg-transparent hover:bg-[#454545] transition-colors duration-200"
                >
                    Submit
                </button>
            </div>
        </Drawer>
    );
};

export default FeedbackDrawer;