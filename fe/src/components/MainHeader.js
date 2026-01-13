import React from "react";
import Image from 'next/image'
import Link from 'next/link'
import bigLogo from '../../public/big.svg';
import appleShortcut from '../../public/shortcut.png';

export default function MainHeader() {

    return (
        <>
            <header className="absolute inset-x-0 top-0 z-50">
                <nav aria-label="Global" className="flex items-center justify-between p-6 lg:px-8">
                    <div className="flex lg:flex-1">
                        <Link href="/" className="-m-1.5 p-1.5">
                            <Image src={bigLogo} alt="logo" width={90} height={90} />
                        </Link>
                    </div>
                    <div className="lg:flex lg:flex-1 lg:justify-end">
                        <Link href="/shorten">
                            <Image src={appleShortcut} alt="apple_shortcut" width={45} height={45} />
                        </Link>
                    </div>
                </nav>
            </header>
        </>
    )
}