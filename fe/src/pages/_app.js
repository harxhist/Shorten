import "@/styles/globals.css";
import { UIProvider } from "@/context/UIContext";
import Loading from "@/components/Loading";
import Notification from "@/components/Notification";
import GoogleAnalytics from "@/components/GoogleAnalytics";

export default function App({ Component, pageProps }) {
  return (
    <>
      {process.env.NODE_ENV === 'production' && (
        <GoogleAnalytics GA_ID={process.env.NEXT_PUBLIC_GA_ID} />
      )}
      <UIProvider>
        <Loading />
        <Notification />
        <Component {...pageProps} />
      </UIProvider>
    </>
  )
}
