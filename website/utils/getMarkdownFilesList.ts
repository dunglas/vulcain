const globby = require('globby');

export async function getMarkdownFilesList(dir: string) {
  const files = await globby([`${dir}/**/*{.md,.mdx}`]);

  return files
    .filter((file) => /\.md$/.test(file))
    .map((file) => ({ params: { slug: file.replace(`${dir}/`, '').replace(/\.md$/, '').split('/') } }));
}
