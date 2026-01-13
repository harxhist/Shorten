import React from 'react';

export default function Footer() {
    return (
        <> 
            <footer className="bg-[#1e1e1e] text-white py-4 mt-auto">
                <div className="flex justify-center space-x-6">
                    <a href={`mailto:${process.env.NEXT_PUBLIC_CONTACT_EMAIL}?subject=[Shorten]`} className="text-yellow-500 hover:text-yellow-400 text-sm">contact</a>
                </div>
            </footer>
        </>
    );

}