import fs from 'fs';
import { join } from 'path';

export function getMarkdown(directory: string, slug: string) {
  const fullPath = join(process.cwd(), directory, `${slug}.md`);
  const fileContents = fs.readFileSync(fullPath, 'utf8');

  return fileContents;
}
