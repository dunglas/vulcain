import React from 'react';
import DocTemplate from '../../components/DocTemplate';
import getAllMarkdownFiles from '../../utils/getGithubMarkdownFiles';
import { GetStaticPaths, GetStaticProps } from 'next';

interface DocPageProps {
  content: string;
}

const DocPage: React.ComponentType<DocPageProps> = ({ content }) => <DocTemplate content={content} />;

export const getStaticPaths: GetStaticPaths = async () => {
  const paths = await getAllMarkdownFiles('docs');

  return {
    paths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const { slug } = params;
  const markdownPath = typeof slug === 'string' ? slug : slug.join('/');
  const response = await fetch(
    `https://raw.githubusercontent.com/${process.env.GITHUB_REPOSITORY}/main/docs/${markdownPath}.md`
  );
  const content = await response.text();
  // Pass data to our component props
  return { props: { content, revalidate: 86400 } };
};
export default DocPage;
