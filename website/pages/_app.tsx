import React, { useEffect } from 'react';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import * as gtag from '../utils/gtag';
import { ThemeProvider } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import theme from '../src/theme';

const App = ({ Component, pageProps }: AppProps) => {
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

  return (
    <>
      <Head>
        <title>Vulcain.rocks: Use HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs</title>
        <link href="https://fonts.googleapis.com/css?family=Oswald:200,300,400,700&display=swap" rel="stylesheet" />
        <link
          href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700,800,900&display=swap"
          rel="stylesheet"
        />
        <meta
          name="description"
          content="Vulcain is a brand new protocol using HTTP/2 Server Push to create fast and idiomatic client-driven REST APIs."
        />
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
