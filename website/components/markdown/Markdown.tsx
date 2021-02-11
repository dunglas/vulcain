import React, { Children } from 'react';
import Link from 'next/link';
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';
import tabs from '../../utils/tabsPlugin';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { Typography, Link as MUILink } from '@material-ui/core';
import { TabPanel } from '@material-ui/lab';
import TabsList from './TabsList';
import CodeBlock from './CodeBlock';
import Heading from './Heading';

const useStyles = makeStyles<Theme>((theme) => ({
  listItem: {
    marginTop: theme.spacing(1),
  },
  link: {
    fontWeight: theme.typography.fontWeightBold,
    color: theme.palette.primary.dark,
  },
  // fit to prism style for inline code.
  inlineCode: {
    color: 'black',
    background: 'rgb(245, 242, 240)',
    fontFamily: 'Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace',
    whiteSpace: 'pre-wrap',
    padding: '5px',
    textShadow: 'white 0px 1px',
  },
  image: {
    maxWidth: '800px',
    margin: '0 auto',
    [theme.breakpoints.down('sm')]: {
      maxWidth: '100%',
    },
  },
}));

interface MarkdownProps {
  source: string;
}

/* eslint-disable react/display-name, react/prop-types */
const Markdown: React.ComponentType<MarkdownProps> = ({ source }) => {
  const classes = useStyles();
  const replacer = (match, p1) => {
    let href;
    if (p1.startsWith('RFC')) {
      href = `https://tools.ietf.org/html/${p1.toLowerCase()}`;
    } else {
      href = `https://duckduckgo.com/\?q=!ducky+site%3Aw3.org+${p1}`;
    }
    return `[@${p1}](${href})`;
  };
  // adds proper href to RFC links
  const formattedSource = source.replace(/\[@!?(.*?)\]/gm, replacer);

  return (
    <ReactMarkdown
      plugins={[gfm, tabs]}
      transformImageUri={(input) => {
        if (/^https?:/.test(input)) return input;
        return `https://raw.githubusercontent.com/dunglas/vulcain/master/${input}`;
      }}
      transformLinkUri={(input) => {
        if (!input || /^#|(https?|mailto):/.test(input)) {
          return input;
        }

        if (input.includes('spec/vulcain.md')) {
          return input.replace(/(.*)#?/, '/spec/vulcain');
        }
        return input.replace('docs/', '/docs/').replace(/\.md/, '');
      }}
      renderers={{
        code: CodeBlock,
        tabs: ({ children, tabs }) => <TabsList tabs={tabs}>{children}</TabsList>,
        tab: ({ value, children }) => <TabPanel value={value}>{children}</TabPanel>,
        inlineCode: ({ value }) => <code className={classes.inlineCode}>{value}</code>,
        paragraph: (props) => <Typography variant="body2" paragraph {...props} />,
        heading: Heading,
        listItem: ({ children }) => {
          // Paragraphs aren't allowed inside list items
          const cleanedChildren = Children.toArray(children).map((child: any) => {
            return child.type && child.type.name === 'paragraph' ? child.props.children : child;
          });

          return (
            <li className={classes.listItem}>
              <Typography component="span">{cleanedChildren}</Typography>
            </li>
          );
        },
        image: (props) => <img className={classes.image} alt={props?.alt} {...props} />,
        link: ({ href, ...props }) => {
          if (!href) {
            return props.children;
          }

          if (/^#|(https?|mailto):/.test(href)) {
            return (
              <MUILink className={classes.link} href={href}>
                {props.children}
              </MUILink>
            );
          }

          return (
            <Link href={`${href}`} passHref>
              <MUILink className={classes.link}>{props.children}</MUILink>
            </Link>
          );
        },
      }}
    >
      {formattedSource}
    </ReactMarkdown>
  );
};
export default Markdown;
