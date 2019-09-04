var resources = {
    "/books": {
        "@id": "/books",
        "hydra:member": [
            "/books/1",
            "/books/2"
        ]
    },
    "/books/1": {
        "@id": "/books/1",
        "title": "Book 1",
        "description": "A good book",
        "author": "/authors/1",
    },
    "/books/2": {
        "@id": "/books/2",
        "title": "Book 2",
        "description": "A great book",
        "author": "/authors/1",
    },
    "/authors/1": {
        "@id": "/authors/1",
        "name": "Author 1",
    }
}

function fixtures(r) {
    r.log(r.uri);

    if (r.uri in resources) {
        r.return(200, JSON.stringify(resources[r.uri]));

        return;
    }

    r.return(404, "Not Found");
}
