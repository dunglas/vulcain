import React from 'react';
import { Drawer, Box, Hidden } from '@material-ui/core';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { ThemeProvider } from '@material-ui/styles';
import docTheme from '../src/docTheme';
import DocsMenu from './DocsMenu';
import Markdown from './markdown/Markdown';
import Page from './Page';

const drawerWidth = 240;

const useStyles = makeStyles<Theme>((theme) => ({
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
    position: 'sticky',
    top: '64px',
    height: 'calc(100vh - 64px)',
  },
  docWrapper: {
    width: `calc(100% - ${drawerWidth}px)`,
    [theme.breakpoints.down('sm')]: {
      width: '100%',
    },
  },
}));

interface DocTemplateProps {
  content: string;
}

const DocTemplate: React.ComponentType<DocTemplateProps> = ({ content }) => {
  const classes = useStyles();

  return (
    <Page withFooter={false}>
      <Box display="flex" width="100%">
        <Hidden smDown initialWidth="md">
          <Drawer
            className={classes.drawer}
            variant="permanent"
            classes={{
              paper: classes.drawerPaper,
            }}
          >
            <DocsMenu />
          </Drawer>
        </Hidden>
        <ThemeProvider theme={docTheme}>
          <Box p={3} className={classes.docWrapper}>
            {content && <Markdown source={content} />}
          </Box>
        </ThemeProvider>
      </Box>
    </Page>
  );
};

export default DocTemplate;
