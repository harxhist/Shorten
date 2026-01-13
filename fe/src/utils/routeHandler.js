// import { fetchSummaryContent } from '../utils/apiHandler';
import { getIronSession } from 'iron-session';
import sessionOptions from '../components/sessionOptions';

export const getServerSidePropsForDynamicRoute = async (context) => {
  function isValidExternalUrl(url) {
    try {

      // Ensure the URL is not pointing to localhost:3000
      if (url.includes('localhost') ||  url.includes('3000')) {
        return false; // Localhost URL is not external
      }

      // Check if the URL contains typical Next.js identifiers
      const nextJsPatterns = [
        '_next', // Next.js assets
        '/api/'  // Next.js API routes
      ];

      for (const pattern of nextJsPatterns) {
        if (url.includes(pattern)) {
          return false; // URL is likely a Next.js internal URL
        }
      }

      // If all checks pass, it's a valid external URL
      return true;
    } catch (e) {
      // If URL parsing fails, it's not a valid URL
      return false;
    }
  }


  const { req, res, params } = context;

  if (!params?.path || params.path.length === 0) {
    return { redirect: { destination: '/', permanent: false } };
  }
  let url;
  if (params.path[0] === 'http:' || params.path[0] === 'https:') {
    url = params.path.slice(1).join('/');
  } else {
    url = decodeURIComponent(params.path.join('/'));
  }
  try {
    if (!isValidExternalUrl(url)) {
      req.session = null;
      return {
        redirect: {
          destination: '/',
          permanent: false,
        },
      };
    }else{
      const session = await getIronSession(req, res, sessionOptions());
    session.url = url;
    await session.save();
    return {
      redirect: {
        destination: '/gist',
        permanent: false,
      },
    };
    }
    
  } catch (error) {
    console.error('Error saving URL to session:', error);
    return { redirect: { destination: '/', permanent: false } };
  }
};