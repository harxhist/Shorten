import { useRouter } from 'next/router';
import { useState, useEffect, useRef } from 'react';
import { getIronSession } from 'iron-session';
import sessionOptions from '../components/sessionOptions';
import { fetchSummaryContent } from '../utils/apiHandler';
import Link from 'next/link';
import Image from 'next/image';
import small from '../../public/small.svg';
import AutoSizingCard from '../components/AutoSizingCard';
import Footer from '../components/Footer';
import { putFeedback } from '../utils/apiHandler';
import { getCachedData } from '@/utils/cacheUtils';
import { Tooltip, IconButton } from '@mui/material';


export const getServerSideProps = async (context) => {
  const { req, res } = context;
  const session = await getIronSession(req, res, sessionOptions());
  const url = session.url || null;

  return {
    props: {
      initialUrl: url,
    },
  };
};

const SumPage = ({ initialUrl }) => {
  const [content, setContent] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const router = useRouter();

  useEffect(() => {
    // Handle route changes
    const handleRouteChange = (url) => {
      // Function to run when user navigates away from /gist page
      if (router.pathname === '/gist') {
        putFeedback();
        // Add your cleanup logic here
      }
    };

    // Handle window/tab close
    const handleBeforeUnload = (e) => {
      putFeedback();
    };

    // Set up event listeners
    router.events.on('routeChangeStart', handleRouteChange);
    window.addEventListener('beforeunload', handleBeforeUnload);

    // Clean up event listeners
    return () => {
      router.events.off('routeChangeStart', handleRouteChange);
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [router]);

  useEffect(() => {
    const fetchData = async () => {
      if (initialUrl) {
        setIsLoading(true);
        setError(null);
        try {
          const key = `summary_${initialUrl}`;
          if (getCachedData(key)) {
            setContent(getCachedData(key));
            setIsLoading(false);
            return;
          } else {
            localStorage.clear();
            const fetchedContent = await fetchSummaryContent(initialUrl);
            setContent(fetchedContent);
            await fetch('/api/clear-session', { method: 'POST' });
          }
        } catch (err) {
          setError(err.message || 'An error occurred');
        } finally {
          setIsLoading(false);
        }
      }
    };

    fetchData();
  }, [initialUrl]);

  return (
    <div className="bg-[#1e1e1e] h-screen overflow-hidden flex flex-col">
      <header className="absolute inset-x-0 top-0 z-50">
        <nav className="flex items-center justify-between p-6 lg:px-8" aria-label="Global">
          <div className="flex lg:flex-1">
            <Link href="/" className="-m-1.5 p-1.5">
              <Image src={small} alt="logo" width={40} height={40} />
            </Link>
          </div>
          
          <div >

        <span
          type="button"
          className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border border-blue-400 bg-blue-950 text-blue-300"
          
        >
          BETA
        </span>
        
        
          
    
    </div>

      
        </nav>
      </header>
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

        <div className="mx-auto max-w-2xl py-0 sm:py-12 lg:py-14 text-center">
          <div className="mb-20 flex justify-center "></div>

          <article className="prose dark:prose-invert prose-img:rounded-xl prose-headings:underline prose-a:text-blue-600">
              <AutoSizingCard mdText={content} loading={isLoading} error={error} />
          </article>

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
      </div>
      <Footer />
    </div>
  );
};

export default SumPage;