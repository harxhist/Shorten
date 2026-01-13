import React from "react";
import Image from 'next/image';
import small from '../../public/small.svg';
import { useRouter } from 'next/router';

export default function Header() {
    const router = useRouter();
    const handleHome = () => {
        router.push('/');
    };

    return (
        <header className="sticky top-0 py-4 z-50 w-full bg-[#1e1e1e]"> 
            <div className="flex justify-center items-center"> 
                <div className="cursor-pointer" onClick={handleHome}> 
                  <Image src={small} alt="logo" width={40} height={40}/>
                </div>
            </div>
        </header>
    );
}