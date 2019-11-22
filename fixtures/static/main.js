const apiURL = "https://localhost:3000";

const cache = {};
const result = document.getElementById('result');
async function fetchRel(rel) {
  // Prevent fetching twice the same relation
  if (cache[rel]) {
    return cache[rel];
  }

  // use a Promise to wait for pushed relation in the local cache
  let res;
  cache[rel] = new Promise((resolve) => { res = resolve });

  const resp = await fetch(apiURL + rel, { credentials: "include" });
  const json = await resp.json();
  res(json);
  return json;
}

(async function() {
  const books = await fetchRel("/books.jsonld?preload=/hydra:member/*/author");
  result.innerText = JSON.stringify(books, null, 2);

  books["hydra:member"].forEach(async (bookId, i) => {
    const book = await fetchRel(bookId);
    book.author = await fetchRel(book.author);
    books["hydra:member"][i] = book;
    result.innerText = JSON.stringify(books, null, 2);
  });

  return books;
})();
