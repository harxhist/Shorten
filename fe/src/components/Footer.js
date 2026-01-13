import React from 'react';
import Link from 'next/link';

export default function Footer() {
    return (
        <> 
            <footer className="bg-[#1e1e1e] text-white py-4 mt-auto">
                <div className="flex flex-col sm:flex-row justify-center items-center space-y-2 sm:space-y-0 sm:space-x-6">
                    <a 
                        href={`mailto:${process.env.NEXT_PUBLIC_CONTACT_EMAIL}?subject=[Shorten]`} 
                        className="text-yellow-500 hover:text-yellow-400 text-sm"
                    >
                        contact
                    </a>
                    <span className="hidden sm:inline text-gray-600">|</span>
                    <Link 
                        href="/privacy" 
                        className="text-yellow-500 hover:text-yellow-400 text-sm"
                    >
                        privacy
                    </Link>
                    <span className="hidden sm:inline text-gray-600">|</span>
                    <Link 
                        href="/terms" 
                        className="text-yellow-500 hover:text-yellow-400 text-sm"
                    >
                        terms
                    </Link>
                </div>
            </footer>
        </>
    );
}