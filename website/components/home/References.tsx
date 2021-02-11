import React from 'react';
import { Typography, Container, Grid, Theme } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { logos } from '../../data/references';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    padding: theme.spacing(8, 0),
    backgroundColor: theme.palette.grey[100],
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
  logoImage: {
    maxWidth: '90%',
    height: 'auto',
  },
}));

const References: React.ComponentType = () => {
  const classes = useStyles();

  return (
    <section className={classes.root}>
      <Container>
        <Typography className={classes.title} align="center" variant="h3" color="textPrimary">
          They use Vulcain
        </Typography>
        <Grid container justify="center">
          {logos.map((logo) => (
            <Grid item xs={4} sm={3} md={2} key={logo.name}>
              <img className={classes.logoImage} src={`static/references/${logo.logo}.png`} alt={logo.name} />
            </Grid>
          ))}
        </Grid>
      </Container>
    </section>
  );
};

export default References;
