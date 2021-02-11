import React from 'react';
import { Box, Typography, Container, Grid, Button, Hidden, Link as MuiLink, Theme } from '@material-ui/core';
import Link from 'next/link';
import { makeStyles } from '@material-ui/core/styles';
import Animation from './Animation';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    background: `url("/static/isometric.png") 140%, linear-gradient(315deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 100%)`,
    color: '#fff',
    padding: theme.spacing(8, 0),
  },
  image: {
    [theme.breakpoints.down('sm')]: {
      paddingBottom: theme.spacing(3),
    },
  },
  button: {
    borderRadius: '40px',
    marginTop: theme.spacing(3),
    '&:not(:first-child)': {
      marginLeft: theme.spacing(1),
    },
  },
  content: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
    },
    [theme.breakpoints.up('md')]: {
      textAlign: 'left',
    },
  },
}));

const Main: React.ComponentType = () => {
  const classes = useStyles();

  return (
    <section className={classes.root}>
      <Container>
        <Grid container alignItems="center" justify="center">
          <Hidden xsDown initialWidth="md">
            <Grid item xs={12} md={7} className={classes.image}>
              <Animation />
            </Grid>
          </Hidden>
          <Grid item md={5}>
            <Box className={classes.content} p={1}>
              <Typography variant="h1" color="inherit">
                Vulcain:
              </Typography>
              <Typography variant="h2" color="inherit" gutterBottom>
                client-driven hypermedia APIs
              </Typography>
              <Typography paragraph>
                Vulcain is a brand new protocol using HTTP/2 Server Push to create fast and idiomatic{' '}
                <strong>client-driven REST</strong> APIs.
              </Typography>
              <Typography paragraph>
                An open source gateway server which you can put on top of <strong>any existing web API</strong> to
                instantly turn it into a Vulcain-compatible one is also provided!
              </Typography>
              <Typography paragraph>
                It supports{' '}
                <Box fontWeight="fontWeightBold" component="span">
                  <MuiLink
                    href="https://restfulapi.net/hateoas/"
                    target="_blank"
                    rel="noopener noreferrer"
                    underline="always"
                    color="inherit"
                  >
                    hypermedia APIs
                  </MuiLink>
                </Box>{' '}
                but also any &quot;legacy&quot; API by documenting its relations{' '}
                <Box fontWeight="fontWeightBold" component="span">
                  <MuiLink
                    href="https://github.com/dunglas/vulcain/blob/master/docs/gateway/openapi.md"
                    target="_blank"
                    rel="noopener noreferrer"
                    underline="always"
                    color="inherit"
                  >
                    using OpenAPI
                  </MuiLink>
                </Box>
                .
              </Typography>
              <Box>
                <Link href="/docs" passHref>
                  <Button className={classes.button} size="large" color="secondary" variant="contained" component="a">
                    Get started!
                  </Button>
                </Link>
              </Box>
            </Box>
          </Grid>
        </Grid>
      </Container>
    </section>
  );
};

export default Main;
