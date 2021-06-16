import React, { memo } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';

interface CodeBlockProps {
  language: string;
  value: string;
}

const CodeBlock: React.ComponentType<CodeBlockProps> = ({ language, value }) => {
  const codeTagProps = {
    style: {
      color: 'black',
      background: 'none',
      textShadow: 'white 0px 1px',
      fontFamily: 'Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace',
      textAlign: 'left',
      whiteSpace: 'pre-wrap',
      wordSpacing: 'normal',
      wordBreak: 'break-all',
      overflowWrap: 'normal',
      lineHeight: '1.5',
      tabSize: '4',
      hyphens: 'none',
    },
  };
  return (
    <SyntaxHighlighter
      language={language}
      customStyle={{ display: 'block', whiteSpace: 'pre-wrap' }}
      codeTagProps={codeTagProps}
    >
      {value}
    </SyntaxHighlighter>
  );
};

export default memo(CodeBlock);
