import React from 'react';
import DocTemplate from '../../components/DocTemplate';
import { getFiles } from '../../utils/getAllFolderFileNames';
import { GetStaticPaths, GetStaticProps } from 'next';
import { getMarkdown } from '../../utils/getMarkdownByFilePath';

interface DocPageProps {
  content: string;
}

const DocPage: React.ComponentType<DocPageProps> = ({ content }) => <DocTemplate content={content} />;

export const getStaticPaths: GetStaticPaths = async () => {
  const paths = await getFiles('docs');
  return {
    paths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const { slug } = params;
  const markdownPath = typeof slug === 'string' ? slug : slug.join('/');
  const content = getMarkdown('docs', markdownPath);
  // Pass data to our component props
  return { props: { content } };
};

export default DocPage;
