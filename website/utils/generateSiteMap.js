const fs = require('fs');
const globby = require('globby');

const ROOT_URL = 'https://vulcain.rocks';

async function generateSiteMap() {
  const pages = await globby([
    'pages/**/*{.ts,.tsx}',
    'docs/**/*.md',
    'spec/**/*.md',
    '!pages/**/[*.{ts,tsx}',
    '!pages/_*.{ts,tsx,.js,.jsx}',
    '!pages/**/[...slug]{.ts,.tsx}',
    '!docs/**/README.md',
  ]);
  console.log(pages);
  const sitemap = `
      <?xml version="1.0" encoding="UTF-8"?>
      <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          ${pages
            .map((page) => {
              const path = page
                .replace('pages/', '')
                .replace('.tsx', '')
                .replace('.md', '')
                .replace('.ts', '')
                .replace('index', '');
              const route = path === '/index' ? '' : path;
              return `
                      <url>
                      <loc>${`${ROOT_URL}/${route}`}</loc>
                      </url>
                  `;
            })
            .join('')}
      </urlset>
  `;

  fs.writeFileSync('public/sitemap.xml', sitemap);
}

generateSiteMap();
