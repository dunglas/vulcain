import React from 'react';
import DocTemplate from '../../components/DocTemplate';
import { getMarkdownFilesList } from '../../utils/getMarkdownFilesList';
import { GetStaticPaths, GetStaticProps } from 'next';
import { getMarkdown } from '../../utils/getMarkdownByFilePath';

interface SpecPageProps {
  content: string;
}

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

const SpecPage: React.ComponentType<SpecPageProps> = ({ content }) => <DocTemplate content={cleanSpec(content)} />;

export const getStaticPaths: GetStaticPaths = async () => {
  const paths = await getMarkdownFilesList('spec');
  return {
    paths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const { slug } = params;
  const markdownPath = typeof slug === 'string' ? slug : slug.join('/');
  const content = getMarkdown('spec', markdownPath);
  // Pass data to our component props
  return { props: { content } };
};
export default SpecPage;
