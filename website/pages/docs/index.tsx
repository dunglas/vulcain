import React from 'react';
import DocTemplate from '../../components/DocTemplate';
import { GetStaticProps } from 'next';
import { getMarkdown } from '../../utils/getMarkdownByFilePath';

interface DocPageProps {
  content: string;
}

const DocPage: React.ComponentType<DocPageProps> = ({ content }) => <DocTemplate content={content} />;

export const getStaticProps: GetStaticProps = async () => {
  const content = getMarkdown('docs', 'README');
  // Pass data to our component props
  return { props: { content } };
};

export default DocPage;
