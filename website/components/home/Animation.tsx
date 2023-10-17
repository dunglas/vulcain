import React, { useEffect, useRef, useState } from 'react';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { useIntersection } from 'react-use';
import { Box, Paper } from '@material-ui/core';
import gsap, { Sine as Cubic } from 'gsap';
import Head from 'next/head';
import MethodSelector from '../MethodSelector';
import { METHODS } from '../../data/methods';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    position: 'relative',
    opacity: 0,
    maxWidth: '600px',
    margin: '0 auto',
    padding: theme.spacing(2),
    borderLeft: `20px solid ${theme.palette.primary.light}`,
    transform: 'rotate(-2deg)',
  },
  base: {
    width: '80%',
    height: 'auto',
    position: 'relative',
  },
  api: {
    position: 'absolute',
    width: '25%',
    top: '55%',
    left: '70%',
    transform: 'rotate(2deg) translateY(-50%)',
    '& img': {
      width: '100%',
      height: 'auto',
    },
    [theme.breakpoints.down('md')]: {
      top: '50%',
      left: '72%',
    },
  },
  animated: {
    position: 'absolute',
    left: 0,
    top: 0,
    width: '80%',
    height: 'auto',
  },
}));

const Animation: React.ComponentType = () => {
  const container = useRef(null);

  const classes = useStyles();
  const [methodKey, setMethodKey] = useState(Object.keys(METHODS)[0]);

  const method = METHODS[methodKey];

  const timeline = gsap.timeline({ paused: true });

  useEffect(() => {
    timeline.set('#step1', {
      clipPath: 'polygon(0% 0%, 0% 0%, 0% 100%, 0% 100%)',
    });
    timeline.set('#step2', {
      clipPath: 'polygon(100% 0%, 100% 0%, 100% 100%, 100% 100%)',
    });
    timeline.set('#step2-base', {
      opacity: 0,
    });
    if (method.steps > 2) {
      timeline.set('#step3', {
        clipPath: 'polygon(0% 0%, 0% 0%, 0% 100%, 0% 100%)',
      });
    }
    if (method.steps > 3) {
      timeline.set('#step4-base', {
        opacity: 0,
      });
      timeline.set('#step4', {
        clipPath: 'polygon(100% 0%, 100% 0%, 100% 100%, 100% 100%)',
      });
    }
    timeline.to('#step1', 1, {
      clipPath: 'polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)',
      ease: Cubic.easeIn,
    });
    timeline.to(
      '#step2-base',
      0.5,
      {
        opacity: 1,
        ease: Cubic.easeOut,
      },
      'step2'
    );
    timeline.to(
      '#step2',
      1,
      {
        clipPath: 'polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)',
        ease: Cubic.easeOut,
      },
      'step2'
    );
    if (method.steps > 2) {
      timeline.to('#step3', 1, {
        clipPath: 'polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)',
        ease: Cubic.easeIn,
        delay: 0.5,
      });
    }
    if (method.steps > 3) {
      timeline.to(
        '#step4-base',
        0.5,
        {
          opacity: 1,
          ease: Cubic.easeOut,
        },
        'step4'
      );
      timeline.to(
        '#step4',
        1,
        {
          clipPath: 'polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)',
          ease: Cubic.easeOut,
        },
        'step4'
      );
    }

    return () => {
      timeline.kill();
    };
  }, [timeline]);

  const intersection = useIntersection(container, {
    root: null,
    threshold: 0.5,
  });

  if (intersection) {
    if (intersection.isIntersecting) {
      gsap.to('#container', 0.5, { opacity: 1, x: 0 });
      gsap.to('#api', 0.5, { opacity: 1, x: 0 });
      timeline.restart();
    } else {
      gsap.to('#container', 0.5, { opacity: 0, x: 30 });
      gsap.to('#api', 0.5, { opacity: 0, x: -30 });
    }
  }

  return (
    <>
      <Head>
        <link rel="preload" as="image" href="/img/main-schema/base.svg" />
        <link rel="preload" as="image" href="/img/main-schema/API.svg" />
      </Head>
      <Box pl={2} pr={6} position="relative">
        <Paper elevation={20} square className={classes.root} ref={container} id="container">
          <Box position="relative">
            <img className={classes.base} src={'/img/main-schema/base.svg'} alt="main" width="800" height="821" />
            <img
              className={classes.animated}
              src={`/img/main-schema/${method.folder}/step1.svg`}
              alt="step1"
              id="step1"
              width="800"
              height="821"
            />
            <img
              className={classes.animated}
              src={`/img/main-schema/${method.folder}/step2.svg`}
              alt="step2"
              id="step2"
              width="800"
              height="821"
            />
            <img
              className={classes.animated}
              src={`/img/main-schema/${method.folder}/step2-base.svg`}
              alt="step2 circle"
              id="step2-base"
              width="800"
              height="821"
            />
            {method.steps > 2 && (
              <img
                className={classes.animated}
                src={`/img/main-schema/${method.folder}/step3.svg`}
                alt="step3"
                id="step3"
                width="800"
                height="821"
              />
            )}
            {method.steps > 3 && (
              <>
                <img
                  className={classes.animated}
                  src={`/img/main-schema/${method.folder}/step4.svg`}
                  alt="step4"
                  id="step4"
                  width="800"
                  height="821"
                />
                <img
                  className={classes.animated}
                  src={`/img/main-schema/${method.folder}/step4-base.svg`}
                  alt="step4 circle"
                  id="step4-base"
                  width="800"
                  height="821"
                />
              </>
            )}
          </Box>
          <MethodSelector method={methodKey} onMethodChange={(method) => setMethodKey(method)} />
        </Paper>
        <Paper elevation={5} square className={classes.api} id="api">
          <img src="/img/main-schema/API.svg" alt="api" width="493" height="904" />
        </Paper>
      </Box>
    </>
  );
};

export default Animation;
