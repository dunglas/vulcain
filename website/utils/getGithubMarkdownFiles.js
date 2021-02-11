require('isomorphic-fetch');

async function getAllFiles(dir) {
  const response = await fetch(
    `https://api.github.com/repos/${process.env.GITHUB_REPOSITORY}/contents/${dir}?ref=main`,
    {
      headers: {
        authorization: `token ${process.env.GITHUB_TOKEN}`,
      },
    }
  );
  const contents = await response.json();
  const files = await Promise.all(
    contents.map((content) => {
      if (content.type === 'dir') {
        console.log('there is a dir', content.path);
        return getAllFiles(content.path);
      } else return content.path;
    })
  );
  return Array.prototype.concat(...files);
}

export default async function getFiles(dir) {
  const files = await getAllFiles(dir);

  return files
    .filter((file) => /\.md$/.test(file))
    .map((file) => ({ params: { slug: file.replace(`${dir}/`, '').replace(/\.md$/, '').split('/') } }));
}
