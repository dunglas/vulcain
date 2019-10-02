const apiURL = "https://localhost:3000";

const cache = {};
async function fetchRel(rel) {
  // Prevent a bug (?) in Chrome
  if (cache[rel]) {
    return cache[rel];
  }

  const resp = await fetch(apiURL + rel, { credentials: "include" });
  const json = await resp.json();
  cache[rel] = json;

  return json;
}

(async function() {
  const books = await fetchRel("/books.jsonld?preload=/hydra:member/*/author");

  for (let i = 0; i < books["hydra:member"].length; i++) {
    books["hydra:member"][i] = await fetchRel(books["hydra:member"][i]);
    books["hydra:member"][i].author = await fetchRel(
      books["hydra:member"][i].author
    );
  }

  document.write(`<pre><code>${JSON.stringify(books, null, 2)}</code></pre>`);

  return books;
})();
