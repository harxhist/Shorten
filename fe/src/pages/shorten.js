import React from 'react';
import Image from 'next/image';
import Footer from '../components/Footer';
import ios from '../../public/shorten.png';
import shortcut from '../../public/shortcut.png';
import Header from '../components/Header';
import one from '../../public/1.png';
import two from '../../public/2.png';
import three from '../../public/3.png';
import four from '../../public/4.png';

const Shorten = () => {
  const steps = [
    { image: one, title: 'Normal Browsing ', description: 'The user visits any blog or webpage & needs TL;DR for it' },
    { image: two, title: 'Adding the Prefix', description: 'Before the original URL, insert "sh10.io/".' },
    { image: three, title: 'Getting Gist', description: 'Summary will appear. Click "Listen" to hear it.' },
  ];

  return (
    <div className="bg-[#1e1e1e] antialiased min-h-screen">
      <Header />
      <div className="bg-[#282828] py-0 px-4 md:px-6">
        <div className="max-w-screen-lg mx-auto flex flex-col md:flex-row items-center md:space-x-12">
          <div className="w-1/2 hidden md:block">
            <Image src={ios} alt="shorten" width={500} height={500} />
          </div>
          <div className="pb-8 text-center md:text-left">
            <a href="/">
            </a>
            <h1 className="mt-12 text-4xl text-white font-thin">Summarise w/ shorten</h1>
            <p className="mt-4 text-base text-white font-light">Get a gist and listen to it.</p>
            <a href="/ios-shortcut" target="_blank">
              <div className="mt-12 inline-flex items-center py-1.5 px-4 rounded-xl space-x-2 text-white hover:shadow-lg bg-[#1D1F57]">
                <Image src={shortcut} width={60} height={60} />
                <div className="px-2">
                  <p className="font-light leading-none">Download the</p>
                  <h3 className="text-3xl font-semibold leading-none">
                    Shortcut
                  </h3>
                </div>
              </div>
            </a>
          </div>
          <div className="md:hidden text-white px-12">
            <Image src={ios} alt="shorten" width={500} height={500} />
          </div>
        </div>
      </div>
      <div className="max-w-screen-lg mx-auto px-4 md:px-6 py-8">
        <div className="py-20 px-4 md:px-6 flex flex-col md:flex-row gap-12 items-center justify-between">

          <div className="flex flex-col space-y-4">
            <h2 className="text-3xl text-white font-thin">Simplify Jargon</h2>
            <ul className="text-white list-none space-y-2">
              <li>Summarizes most webpages.</li>
              <li>Listen to summarized webpage.</li>
            </ul>
          </div>
          <div className="flex flex-col space-y-4">
            <h2 className="text-3xl text-white font-thin"> Simple Ways</h2>
            <ul className="text-white list-none space-y-2">
              <li>Prefix any URL with "sh10.io/"</li>
              <li>Use <em>Shorten</em> shortcut on Mac/ios.</li>
            </ul>
          </div>
        </div>
        <div className="prose-lg text-white gap-6">
          <h2 className='md:px-6 px-4 text-3xl text-white font-normal text-center mt-10'>Guide</h2>
          <br />
          <br />
          <div className="space-y-12">
            {steps.map((step, index) => (
              <div key={index} className={`md:flex items-center ${index % 2 === 0 ? 'md:flex-row' : 'md:flex-row-reverse'} md:space-x-6`}>
                <div className="md:w-1/2 relative p-6">
                  <div className="md:w-[70%] md:h-[70%] w-[60%] h-[60%] mx-auto">
                    <Image
                      src={step.image}
                      width={500}
                      height={300}
                      alt={step.title}
                      className="object-cover w-full h-full block"
                    />
                  </div>
                </div>
                <div className="md:w-1/2 p-6">
                  <div className="md:block hidden">
                    <h3 className='text-[#f2bb4b] text-lg'>{step.title}</h3>
                    <p>{step.description}</p>
                  </div>
                  <div className="md:hidden text-center">
                    <h3 className='text-[#f2bb4b] text-lg'>{step.title}</h3>
                    <p>{step.description}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Shorten;