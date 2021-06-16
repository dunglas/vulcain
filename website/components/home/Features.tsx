import React from 'react';
import { Typography, Container, Grid, Theme } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import {
  LockOutlined,
  SecurityOutlined,
  Http,
  ImportantDevicesOutlined,
  SendOutlined,
  CompareArrowsOutlined,
  ArchiveOutlined,
  LanguageOutlined,
  MessageOutlined,
  People,
  Storage,
  Speed,
} from '@material-ui/icons';
import useAnimation, { DirectionType } from '../../hooks/useAnimation';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    padding: theme.spacing(8, 0),
    backgroundColor: theme.palette.grey[100],
    overflow: 'hidden',
  },
  title: {
    marginBottom: theme.spacing(10),
    position: 'relative',
    '&::after': {
      content: "''",
      position: 'absolute',
      width: '100px',
      height: '8px',
      backgroundColor: theme.palette.secondary.main,
      bottom: '-20px',
      left: '50%',
      transform: 'translateX(-50%)',
    },
  },
  button: {
    borderRadius: '40px',
    marginTop: theme.spacing(3),
  },
  featureItem: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  featureIcon: {
    fontSize: '2.2rem',
  },
  featureTitle: {
    maxWidth: '250px',
  },
  featureIconCircle: {
    backgroundColor: theme.palette.primary.dark,
    borderRadius: '50%',
    width: '70px',
    height: '70px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: '#fff',
    position: 'relative',
    marginBottom: theme.spacing(1),
    border: `5px solid ${theme.palette.primary.light}`,
    '&::after': {
      content: "''",
      position: 'absolute',
      width: '60px',
      height: '60px',
      borderRadius: '50%',
      border: `4px solid #fff`,
      left: '50%',
      top: '50%',
      transform: 'translate(-50%, -50%)',
    },
  },
}));

interface FeatureItemProps {
  title: string;
  Icon: React.ElementType;
  classes: any;
  direction?: DirectionType;
}

const FeatureItem: React.ComponentType<FeatureItemProps> = ({ title, Icon, classes, direction }) => {
  const animation = useAnimation(direction, { rootMargin: '-10%' });
  return (
    <Grid item xs={12} sm={4} md={3} lg={2} className={classes.featureItem} ref={animation}>
      <div className={classes.featureIconCircle}>
        <Icon className={classes.featureIcon} />
      </div>
      <Typography className={classes.featureTitle} align="center" variant="subtitle2" component="span" gutterBottom>
        {title}
      </Typography>
    </Grid>
  );
};

const Features: React.ComponentType = () => {
  const classes = useStyles();

  return (
    <section className={classes.root}>
      <Container>
        <Typography className={classes.title} align="center" variant="h3" color="primary">
          Vulcain: at a glance
        </Typography>
        <Grid container spacing={2} alignItems="flex-start" justify="center">
          <FeatureItem
            classes={classes}
            Icon={Http}
            title="Pure HTTP, full-duplex, leverage HTTP/2+"
            direction="left"
          />
          <FeatureItem classes={classes} Icon={Speed} title="High performance, low latency" direction="right" />
          <FeatureItem
            classes={classes}
            Icon={ImportantDevicesOutlined}
            title="Native browser support, works everywhere"
            direction="left"
          />
          <FeatureItem
            classes={classes}
            Icon={SendOutlined}
            title="Publish with a simple POST request"
            direction="right"
          />
          <FeatureItem
            classes={classes}
            Icon={CompareArrowsOutlined}
            title="Subscribe using Server-Sent-Events"
            direction="left"
          />
          <FeatureItem
            classes={classes}
            Icon={ArchiveOutlined}
            title="Automatic reconnection, refetch missed messages"
            direction="right"
          />
          <FeatureItem
            classes={classes}
            Icon={MessageOutlined}
            title="Designed for REST and GraphQL"
            direction="left"
          />
          <FeatureItem
            classes={classes}
            Icon={SecurityOutlined}
            title="Private updates (JWT authorization)"
            direction="right"
          />
          <FeatureItem classes={classes} Icon={People} title="Presence API and subscription events" direction="left" />
          <FeatureItem classes={classes} Icon={Storage} title="Event store" direction="right" />
          <FeatureItem
            classes={classes}
            Icon={LanguageOutlined}
            title="Compatible with serverless, PHP and the like"
            direction="left"
          />
          <FeatureItem classes={classes} Icon={LockOutlined} title="Supports end-to-end encryption" direction="right" />
        </Grid>
      </Container>
    </section>
  );
};

export default Features;
