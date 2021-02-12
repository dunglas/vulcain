// see https://github.com/rexxars/react-markdown/issues/69
import React from 'react';
import Typography from '@material-ui/core/Typography';
import Head from 'next/head';

function flatten(text, child) {
  return typeof child === 'string' ? text + child : React.Children.toArray(child.props.children).reduce(flatten, text);
}

interface HeadingProps {
  level: 1 | 2 | 3 | 4 | 5;
}

type TitleVariantType = 'h1' | 'h2' | 'h3' | 'h4' | 'h5';

const Heading: React.ComponentType<HeadingProps> = ({ children, level }) => {
  const text = React.Children.toArray(children).reduce(flatten, '');

  let slug, found;
  if ((found = text.match(/{#(.*)}/))) {
    slug = found[1];
    children = text.replace(/ {#.*}/, '');
  } else slug = text.toLowerCase().replace(/\W/g, '-');

  // TODO: Create clickable links to anchors, something like <Link href={`#${slug}`}><LinkIcon /></Link>
  const t = (
    <Typography variant={`h${level}` as TitleVariantType} id={slug}>
      {children}
    </Typography>
  );

  if (level !== 1) {
    return t;
  }

  const title = text;
  const schema = {
    '@context': 'https://schema.org/',
    '@type': 'TechArticle',
    name: text,
  };

  // create meta tags thanks to first title
  return (
    <>
      <Head>
        <title>{title} - Vulcain.rocks</title>
        <meta name="description" content={title} />
        <meta name="og:title" content={`Vulcain.rocks - ${title}`} />
        <meta name="twitter:title" content={`Vulcain.rocks - ${title}`} />
        <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(schema) }} />
      </Head>
      {t}
    </>
  );
};

export default Heading;
