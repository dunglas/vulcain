// Adapted from https://github.com/mui-org/material-ui/tree/master/docs/src/pages/getting-started/page-layout-examples/pricing
// TODO: could be a HOC
import React from 'react';
import classNames from 'classnames';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import { makeStyles, Theme } from '@material-ui/core';
import Header from './Header';
import Link from 'next/link';
import MUILink from '@material-ui/core/Link';

const useStyles = makeStyles<Theme>((theme) => ({
  '@global': {
    body: {
      backgroundColor: theme.palette.common.white,
    },
    '#__next': {
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
    },
  },
  layout: {
    width: 'auto',
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  footer: {
    backgroundColor: theme.palette.grey[900],
    color: '#fff',
    padding: theme.spacing(6, 0),
  },
  footerLink: {
    textDecoration: 'none',
    display: 'block',
  },
  content: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
  },
}));

const footers = [
  {
    title: 'Vulcain',
    content: {
      'Get Started': '/docs',
      Gateway: '/docs/gateway/install',
    },
  },
  {
    title: 'Social',
    content: {
      GitHub: 'https://github.com/dunglas/vulcain',
      'Contact us': 'mailto:dunglas+vulcain@gmail.com',
    },
  },
  {
    title: 'Support',
    content: {
      'Stack Overflow': '/docs/help',
      'Community Slack': '/docs/help',
      Training: 'https://masterclass.les-tilleuls.coop/en/trainings/discover-vulcain',
    },
  },
  {
    title: 'API Platform ecosystem',
    content: {
      'API Platform': 'https://api-platform.com/',
      'Mercure.rocks': 'https://mercure.rocks/',
    },
  },
];

// TODO: transform in HOC
const Page: React.ComponentType<{ withFooter?: boolean }> = ({ children, withFooter = true }) => {
  const classes = useStyles();

  return (
    <>
      <CssBaseline />
      <Header />
      <main className={classes.layout}>
        <div className={classes.content}>{children}</div>
        {/* Footer */}
        {withFooter && (
          <footer className={classNames(classes.footer)}>
            <Container>
              <Grid container spacing={2} justify="space-evenly">
                {footers.map((footer) => (
                  <Grid item xs key={footer.title}>
                    <Typography variant="h6" color="inherit" gutterBottom>
                      {footer.title}
                    </Typography>
                    {Object.entries(footer.content).map(([k, v]) => {
                      if (/^(https?|mailto):/.test(v)) {
                        return (
                          <MUILink href={v} key={k} variant="subtitle1" color="inherit" className={classes.footerLink}>
                            {k}
                          </MUILink>
                        );
                      }

                      return (
                        <Link href={v} key={k} passHref>
                          <Typography variant="subtitle1" color="inherit" component="a" className={classes.footerLink}>
                            {k}
                          </Typography>
                        </Link>
                      );
                    })}
                  </Grid>
                ))}
              </Grid>
            </Container>
          </footer>
        )}
        {/* End footer */}
      </main>
    </>
  );
};

export default Page;
