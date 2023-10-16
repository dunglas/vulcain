// see https://github.com/rexxars/react-markdown/issues/69
import React from 'react';
import PropTypes from 'prop-types';
import Typography from '@material-ui/core/Typography';
import Head from 'next/head';
import { Variant } from '@material-ui/core/styles/createTypography';

function flatten(text, child) {
  return typeof child === 'string' ? text + child : React.Children.toArray(child.props.children).reduce(flatten, text);
}

const Heading = ({ children, level }) => {
  const text = React.Children.toArray(children).reduce(flatten, '');

  let slug, found;
  if ((found = text.match(/{#(.*)}/))) {
    slug = found[1];
    children = text.replace(/ {#.*}/, '');
  } else slug = text.toLowerCase().replace(/\W/g, '-');

  // TODO: Create clickable links to anchors, something like <Link href={`#${slug}`}><LinkIcon /></Link>
  const t = (
    <Typography variant={`h${level}` as Variant} id={slug}>
      {children}
    </Typography>
  );

  if (level !== 1) {
    return t;
  }

  const title = t.props.children[0].props.children;
  return (
    <React.Fragment>
      <Head>
        <title>{title} - Mercure.rocks</title>
        <meta name="description" content={title} />
        <meta name="og:title" content={`Mercure.rocks: ${title}`} />
      </Head>
      {t}
    </React.Fragment>
  );
};

Heading.propTypes = {
  children: PropTypes.oneOfType([PropTypes.arrayOf(PropTypes.node), PropTypes.node]).isRequired,
  level: PropTypes.number.isRequired,
};

export default Heading;
