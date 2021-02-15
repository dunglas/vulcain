import React, { useEffect } from 'react';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import * as gtag from '../utils/gtag';
import { ThemeProvider } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import theme from '../src/theme';
import '../styles/fonts.css';

const App: React.ComponentType<AppProps> = ({ Component, pageProps }) => {
  const router = useRouter();

  useEffect(() => {
    // Remove the server-side injected CSS.
    const jssStyles = document.querySelector('#jss-server-side');
    if (jssStyles) {
      jssStyles.parentElement.removeChild(jssStyles);
    }
  }, []);

  useEffect(() => {
    const handleRouteChange = (url: URL) => {
      gtag.pageview(url);
    };
    router.events.on('routeChangeComplete', handleRouteChange);
    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router.events]);

  const websiteSchema = {
    '@type': 'WebSite',
    name: 'Vulcain',
    url: 'https://vulcain.rocks',
  };

  return (
    <>
      <Head>
        <title>Vulcain.rocks: Use HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs</title>
        <link rel="preload" href="/fonts/roboto/roboto-v20-latin-regular.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/roboto/roboto-v20-latin-500.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/roboto/roboto-v20-latin-700.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/montserrat/montserrat-v15-latin-200.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/montserrat/montserrat-v15-latin-600.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/montserrat/montserrat-v15-latin-800.woff2" as="font" crossOrigin="" />
        <link rel="preload" href="/fonts/montserrat/montserrat-v15-latin-regular.woff2" as="font" crossOrigin="" />
        <meta
          name="description"
          content="Vulcain is a brand new protocol using HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs."
        />
        <meta name="application-name" content="Vulcain"></meta>
        <meta name="theme-color" content="#f5731b" />
        <meta property="og:url" content="https://vulcain.rocks" />
        <meta property="og:type" content="website" />
        <meta property="og:title" content="Vulcain.rocks" />
        <meta
          property="og:description"
          content="Use HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs"
        />
        <meta property="og:image" content="https://vulcain.rocks/opengraph.png" />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:creator" content="@dunglas" />
        <meta name="twitter:title" content="Vulcain.rocks" />
        <meta
          name="twitter:description"
          content="Use HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs"
        />
        <meta name="twitter:image" content="https://vulcain.rocks/opengraph.png" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="icon" href="/icon.svg" type="image/svg+xml" />
        <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
        <link rel="manifest" href="/site.webmanifest"></link>

        <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(websiteSchema) }} />
      </Head>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        <Component {...pageProps} />
      </ThemeProvider>
    </>
  );
};

export default App;
