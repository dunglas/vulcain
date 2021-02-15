import React from 'react';
import { Box, Typography, Container, Grid, Button, Card, CardContent, CardActions, Theme } from '@material-ui/core';
import Link from 'next/link';
import { KeyboardArrowRight } from '@material-ui/icons';
import SlackIcon from '../icons/Slack';
import StackOverflowIcon from '../icons/StackOverflow';
import { makeStyles } from '@material-ui/core/styles';
import useAnimation from '../../hooks/useAnimation';

const useStyles = makeStyles<Theme>((theme) => ({
  title: {
    marginBottom: theme.spacing(6),
    textAlign: 'left',
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
    },
  },
  root: {
    background: `${theme.palette.primary.dark}`,
    color: '#fff',
    padding: theme.spacing(8, 0),
    position: 'relative',
  },
  container: {
    zIndex: 1,
    position: 'relative',
  },
  cardImage: {
    width: '100%',
    height: 'auto',
  },
  cardCircle: {
    width: '100px',
    height: '100px',
    borderRadius: '50%',
    padding: theme.spacing(2),
    display: 'flex',
    border: `8px solid ${theme.palette.primary.dark}`,
  },
  image: {
    background: `url("/static/help.svg")`,
    backgroundPosition: '0 center',
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'auto 100%',
    [theme.breakpoints.down('sm')]: {
      background: 'none',
    },
  },
  cardContent: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
    height: '100%',
    alignItems: 'flex-start',
    textAlign: 'left',
    justifyContent: 'center',
    [theme.breakpoints.down('sm')]: {
      alignItems: 'center',
      textAlign: 'center',
    },
  },
  cardButton: {
    fontSize: '3rem',
  },
  card: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  cardMain: {
    display: 'flex',
    flex: 1,
    flexDirection: 'row',
    [theme.breakpoints.down('sm')]: {
      flexDirection: 'column',
      textAlign: 'center',
    },
  },
  cardActions: {
    display: 'flex',
    justifyContent: 'flex-end',
    padding: theme.spacing(2),
  },
  svg: {
    position: 'absolute',
    width: '100%',
    height: '30vh',
    zIndex: 0,
    bottom: 0,
    '& > polygon': {
      fill: theme.palette.grey[900],
    },
  },
}));

interface SupportCardProps {
  image: string;
  title: string;
  description: string;
  classes: any;
  children: React.ReactNode;
}

const SupportCard: React.ComponentType<SupportCardProps> = ({ image, title, description, classes, children }) => {
  const animation = useAnimation('bottom', { rootMargin: '-10%' });

  return (
    <Grid item xs={12} sm={6} ref={animation}>
      <Card elevation={3} className={classes.card} square>
        <CardContent className={classes.cardMain}>
          <Box px={2} display="flex" alignItems="center" justifyContent="center">
            <div className={classes.cardCircle}>
              <img src={image} alt={title} className={classes.cardImage} width="150" height="150" />
            </div>
          </Box>
          <Box py={2} className={classes.cardContent}>
            <Typography variant="h5" component="span" gutterBottom className={classes.cardTitle}>
              {title}
            </Typography>
            <Typography variant="body2" color="textSecondary">
              {description}
            </Typography>
          </Box>
        </CardContent>
        <CardActions className={classes.cardActions}>{children}</CardActions>
      </Card>
    </Grid>
  );
};

const Support: React.ComponentType = () => {
  const classes = useStyles();

  return (
    <section className={classes.root}>
      <svg className={classes.svg} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" preserveAspectRatio="none">
        <polygon fill="white" points="0,100 100,0 100,100" />
      </svg>
      <Container className={classes.container}>
        <Grid container alignItems="center" justify="center" className={classes.image}>
          <Grid item xs={false} md={3}></Grid>
          <Grid item xs={12} md={9}>
            <Box className={classes.content} p={1}>
              <Typography variant="h3" color="inherit" className={classes.title}>
                Need help ?
              </Typography>
              <Grid container spacing={4}>
                <SupportCard
                  classes={classes}
                  title="Documentation"
                  description="Reading the documentation is an excellent way to discover Vulcain."
                  image="/static/book.svg"
                >
                  <Link href="/docs" passHref>
                    <Button size="small" color="primary" variant="outlined" component="a">
                      Read the docs
                      <KeyboardArrowRight />
                    </Button>
                  </Link>
                </SupportCard>

                <SupportCard
                  classes={classes}
                  title="Community support"
                  description="Chat with the community on Slack and Stack Overflow"
                  image="/static/handshake.svg"
                >
                  <Link href="/docs/help" passHref>
                    <Button size="small" color="primary" variant="outlined" component="a">
                      Slack
                      <SlackIcon />
                    </Button>
                  </Link>

                  <Link href="/docs/help" passHref>
                    <Button size="small" color="primary" variant="outlined" component="a">
                      Stack Overflow
                      <StackOverflowIcon />
                    </Button>
                  </Link>
                </SupportCard>

                <SupportCard
                  classes={classes}
                  title="Training"
                  description="Improve your Vulcain skills thanks to our trainings."
                  image="/static/presentation.svg"
                >
                  <Button
                    size="small"
                    color="primary"
                    variant="outlined"
                    component="a"
                    href="https://masterclass.les-tilleuls.coop/en/trainings/discover-vulcain"
                  >
                    Our trainings
                    <KeyboardArrowRight />
                  </Button>
                </SupportCard>

                <SupportCard
                  classes={classes}
                  title="Professional services"
                  description="Les-Tilleuls.coop provides professional services: web development, trainings or consulting."
                  image="/static/laptop-code.svg"
                >
                  <Button
                    size="small"
                    color="primary"
                    variant="outlined"
                    component="a"
                    href="https://les-tilleuls.coop"
                  >
                    Les-Tilleuls.coop
                    <KeyboardArrowRight />
                  </Button>
                </SupportCard>
              </Grid>
            </Box>
          </Grid>
        </Grid>
      </Container>
    </section>
  );
};

export default Support;
