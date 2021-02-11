import React, { useEffect, useRef, useState } from 'react';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { useIntersection } from 'react-use';
import { Box, Paper } from '@material-ui/core';
import gsap, { Sine as Cubic } from 'gsap';
import MethodSelector from '../MethodSelector';
import { METHODS } from '../../data/methods';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    position: 'relative',
    opacity: 0,
    maxWidth: '700px',
    margin: '0 auto',
    padding: theme.spacing(2),
    borderLeft: `20px solid ${theme.palette.primary.light}`,
    transform: 'rotate(-2deg)',
  },
  base: {
    width: '80%',
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
    },
  },
  animated: {
    position: 'absolute',
    left: 0,
    top: 0,
    width: '80%',
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

  useEffect(() => {
    gsap.set('#step1', {
      clipPath: 'polygon(0% 0%, 0% 0%, 0% 100%, 0% 100%)',
    });
    gsap.set('#api', { opacity: 0 });
    gsap.set('#result', {
      clipPath: 'polygon(80% 0%, 80% 0%, 80% 100%, 80% 100%)',
    });
  }, []);

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
    <Box pl={2} pr={6} position="relative">
      <Paper elevation={20} square className={classes.root} ref={container} id="container">
        <Box position="relative">
          <img className={classes.base} src={`/static/main-schema/base.png`} alt="main" />
          <img
            className={classes.animated}
            src={`/static/main-schema/${method.folder}/step1.png`}
            alt="step1"
            id="step1"
          />
          <img
            className={classes.animated}
            src={`/static/main-schema/${method.folder}/step2.png`}
            alt="step2"
            id="step2"
          />
          <img
            className={classes.animated}
            src={`/static/main-schema/${method.folder}/step2-base.png`}
            alt="step2 circle"
            id="step2-base"
          />
          {method.steps > 2 && (
            <img
              className={classes.animated}
              src={`/static/main-schema/${method.folder}/step3.png`}
              alt="step3"
              id="step3"
            />
          )}
          {method.steps > 3 && (
            <>
              <img
                className={classes.animated}
                src={`/static/main-schema/${method.folder}/step4.png`}
                alt="step4"
                id="step4"
              />
              <img
                className={classes.animated}
                src={`/static/main-schema/${method.folder}/step4-base.png`}
                alt="step4 circle"
                id="step4-base"
              />
            </>
          )}
        </Box>
        <MethodSelector method={methodKey} onMethodChange={(method) => setMethodKey(method)} />
      </Paper>
      <Paper elevation={5} square className={classes.api} id="api">
        <img src="/static/main-schema/api.png" alt="api" />
      </Paper>
    </Box>
  );
};

export default Animation;
