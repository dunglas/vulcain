import React from 'react';
import DocTemplate from '../../components/DocTemplate';
import { GetStaticProps } from 'next';

interface DocPageProps {
  content: string;
}

const DocPage: React.ComponentType<DocPageProps> = ({ content }) => <DocTemplate content={content} />;

export const getStaticProps: GetStaticProps = async () => {
  const response = await fetch(`https://raw.githubusercontent.com/${process.env.GITHUB_REPOSITORY}/main/README.md`);
  const content = await response.text();
  // Pass data to our component props
  return { props: { content, revalidate: 86400 } };
};
export default DocPage;
