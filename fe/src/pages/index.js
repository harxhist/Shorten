'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import { getServerSidePropsForDynamicRoute } from '../utils/routeHandler'
import MainHeader from '@/components/MainHeader'
import Footer from '../components/Footer'
import { useContext } from 'react'
import { useUI } from '@/context/UIContext'

const dataText = [
  "Short, simple and free. Choose any three.",
];

function typeWriter(text, i, setText, fnCallback) {
  if (i < text.length) {
    setText(text.substring(0, i + 1) + '<span aria-hidden="true"></span>');
    setTimeout(function () {
      typeWriter(text, i + 1, setText, fnCallback)
    }, 40);
  }
  else if (typeof fnCallback == 'function') {
    setTimeout(fnCallback, 100);
  }
}

function StartTextAnimation(i, setText) {
  if (typeof dataText[i] === 'undefined') {
    setTimeout(function () {
      StartTextAnimation(0, setText);
    }, 90000);
  } else {

    typeWriter(dataText[i], 0, setText, function () {
      StartTextAnimation(i + 1, setText);
    });
  }
}

export async function getServerSideProps(context) {
  const { resolvedUrl, query } = context;

  if (query.url) {
    const path = query.url;
    return getServerSidePropsForDynamicRoute({ ...context, params: { path: [path] } });
  }

  return { props: { isHome: true } };
}

export default function Home({ isHome }) {
  const [url, setUrl] = useState('');
  const [text, setText] = useState('');
  const router = useRouter();


  const handleSubmit = (e) => {
    e.preventDefault();
    // showLoading();
    router.push(`/?url=${encodeURIComponent(url)}`);
  };

  useEffect(() => {
    StartTextAnimation(0, setText);
  }, []);

  return (
    <div className="bg-[#1e1e1e] h-screen overflow-hidden flex flex-col">
      <MainHeader/>
      <div className="relative isolate px-6 pt-0 lg:px-8 h-full flex flex-col justify-center items-center">
        <div
          aria-hidden="true"
          className="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
        >
          <div
            style={{
              clipPath:
                'polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)',
            }}
            className="relative left-[calc(50%-11rem)] top-[-3rem] aspect-[800/450] w-[28rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#00c6ff] to-[#0072ff] opacity-50 sm:left-[calc(50%-30rem)] sm:w-[56rem]"
          />
        </div>

        <div className="mx-auto max-w-2xl py-0 sm:py-12 lg:py-14">
          <div className="text-center">
            <h1 className="text-[#f2bb4b] text-5xl font-semibold tracking-tight sm:text-7xl">
              Summarize any webpage.
            </h1>
            <p className="mt-8 text-[#f2bb4b] text-lg font-medium sm:text-xl/8 pb-10">
              <span className="typewriter" dangerouslySetInnerHTML={{ __html: text }} />
            </p>
            <div className="flex justify-center ">
              <form onSubmit={handleSubmit} className="flex text-sm">
                <input
                  className="text-center sm:text-left sm:w-[300px] px-4 w-[150px] border border-gray-600 border-r-0 bg-gray-800 text-white rounded-l-lg focus:outline-none"
                  placeholder="https://gibberish.com/page..."
                  type="text"
                  name="q"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  style={{
                    color: '#f2bb4b',
                    caretColor: '#f2bb4b'
                  }}
                />
                <button
                  type="submit"
                  className="sm:px-4 px-1 py-2 w-full sm:w-auto sm:min-w-6 sm:text-sm text-xs leading-none sm:font-medium border border-yellow-500 bg-yellow-700 text-white rounded-r-lg hover:bg-yellow-600 focus:outline-none"
                >
                  Summarize <br />
                </button>
              </form>
            </div>
            <p className="pt-10 mt-10 text-xs sm:text-xs md:text-sm lg:text-base xl:text-lg text-white">Simply add prefix,</p>
            <p className=" mt-5 text-sm font-style: italic sm:text-base md:text-lg text-white lg:text-xl xl:text-2xl">
              sh10.in/
              <span className="relative inline-block ">
                &lt;URL&gt;
                <span className="convex-underline"></span>
              </span>
            </p>
          </div>
        </div>

        <div
          aria-hidden="true"
          className="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
        >
          <div
            style={{
              clipPath:
                'polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)',
            }}
            className="relative left-[calc(50%+3rem)] top-[-3rem] aspect-[800/450] w-[28rem] -translate-x-1/2 bg-gradient-to-tr from-[#f1c40f] to-[#3498db] opacity-50 sm:left-[calc(50%+36rem)] sm:w-[56rem]"
          />
        </div>
      </div>
      <Footer />
    </div>
  )
}