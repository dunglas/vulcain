import React from 'react';
import { Drawer, Box, Hidden } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { ThemeProvider } from '@material-ui/styles';
import docTheme from '../../src/docTheme';
import DocsMenu from '../../components/DocsMenu';
import Markdown from '../../components/markdown/Markdown';
import Page from '../../components/Page';
import getAllMarkdownFiles from '../../utils/getGithubMarkdownFiles';
import { GetStaticPaths, GetStaticProps } from 'next';

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
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

function cleanSpec(md: string) {
  return (
    '# Vulcain: The Specification\n' +
    md
      .replace(/^%%%(.|\n)*%%%(\W|\n)*\./, '')
      .replace('{mainmatter}', '')
      .replace('{backmatter}', '')
      .replace(/^#\W+/gm, '## ')
      .replace(/\(#(.*?)\)/gm, (_, id) => {
        return `[${id.replace(/-/gm, ' ')}](#${id})`;
      })
  );
}

interface SpecTemplateProps {
  content: string;
}

const SpecTemplate: React.ComponentType<SpecTemplateProps> = ({ content }) => {
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
            <Markdown source={content} />
          </Box>
        </ThemeProvider>
      </Box>
    </Page>
  );
};

export const getStaticPaths: GetStaticPaths = async () => {
  const paths = await getAllMarkdownFiles('spec');

  return {
    paths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const { slug } = params;
  const markdownPath = typeof slug === 'string' ? slug : slug.join('/');
  const response = await fetch(`https://raw.githubusercontent.com/dunglas/vulcain/main/spec/${markdownPath}.md`);
  const content = await response.text();
  // Pass data to our component props
  return { props: { content: cleanSpec(content) } };
};

export default SpecTemplate;
