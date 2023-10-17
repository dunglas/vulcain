/* eslint-disable @typescript-eslint/ban-ts-comment */
import React, { Children } from 'react';
import Link from 'next/link';
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';
import tabs from '../../utils/tabsPlugin';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Link as MUILink } from '@material-ui/core';
import { TabPanel } from '@material-ui/lab';
import TabsList from './TabsList';
import CodeBlock from './CodeBlock';
import Heading from './Heading';
import rehypeRaw from 'rehype-raw';

const useStyles = makeStyles((theme) => ({
  listItem: {
    marginTop: theme.spacing(1),
  },
  link: {
    fontWeight: 'bold',
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
  img: {
    maxWidth: '700px',
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
  const transformImageUri = (input) => {
    if (input.includes('/img/schemas')) {
      return input;
    }
    if (/^https?:/.test(input)) {
      return input;
    }
    if (/^schemas\/vulcain_doc/.test(input)) {
      const result = input.replace('schemas', '/img/schemas');
      return result;
    }
    return `https://raw.githubusercontent.com/dunglas/vulcain/master/${input}`;
  };

  return (
    <ReactMarkdown
      remarkPlugins={[gfm, tabs]}
      // @ts-ignore
      rehypePlugins={[rehypeRaw]}
      transformImageUri={transformImageUri}
      transformLinkUri={(input) => {
        if (!input || /^#|(https?|mailto):/.test(input)) {
          return input;
        }

        if (input.includes('spec/vulcain.md')) {
          return input.replace(/(.*)#?/, '/spec/vulcain');
        }
        return input.replace('docs/', '/docs/').replace(/\.md/, '');
      }}
      components={{
        // @ts-ignore
        code({ inline, className, children, ...props }) {
          const match = /language-(\w+)/.exec(className || '');
          if (className === 'language-tabs') {
            // @ts-ignore
            const tabs = JSON.parse(children);
            return (
              <TabsList tabs={tabs.map((t) => t.value)}>
                {tabs.map((t) => (
                  <TabPanel key={t.value} value={t.value}>
                    <img
                      alt={t.value || ''}
                      className={classes.img}
                      src={transformImageUri(t.children[0].children[0].url)}
                    />
                  </TabPanel>
                ))}
              </TabsList>
            );
          }
          return !inline ? (
            <CodeBlock value={children.toString().trim()} language={match && match[1]} />
          ) : (
            <code className={classes.inlineCode} {...props}>
              {children}
            </code>
          );
        },
        h1: ({ children }) => <Typography variant="h1">{children}</Typography>,
        h2: Heading,
        p: (props) => {
          // @ts-ignore
          return <Typography variant="body2" paragraph {...props} />;
        },
        li: ({ children }) => {
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
        img: (props) => <img className={classes.img} alt={props?.alt} loading="lazy" {...props} />,
        a: ({ href, ...props }) => {
          if (!href) {
            return props.children || '';
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
