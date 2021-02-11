import Page from '../Page';
import { ThemeProvider } from '@material-ui/styles';
import docTheme from '../../src/docTheme';
import { Box } from '@material-ui/core';
import Markdown from './Markdown';
import React from 'react';

const MarkdownPage: React.ComponentType<{ md: string }> = ({ md }) => {
  return (
    <Page>
      <ThemeProvider theme={docTheme}>
        <Box p={3}>
          <Markdown source={md} />
        </Box>
      </ThemeProvider>
    </Page>
  );
};
export default MarkdownPage;
