# Contributing

## Gateway

Clone the project:

    $ git clone https://github.com/dunglas/vulcain
    $ cd vulcain

Install the dependencies:

    $ go get

Run the server:

    $ go run main.go

Run the fixture API:

    $ cd fixtures
    $ go run main.go

Go to `https://localhost:3000` and enjoy!

To run the test suite:

    $ go test -v -timeout 30s github.com/dunglas/vulcain/gateway

### curl Examples

Preload all relations referenced in the `hydra:member`, then in the author relationship, but only include the title and the author of these relations:

```
curl https://localhost:3000/books.jsonld \
    --get \
    --data 'preload=/hydra:member/*/author' \
    --data 'fields=/hydra:member/*/author' \
    --data 'fields=/hydra:member/*/title' \
    --verbose \
    --insecure 
```

Using headers:

```
curl https://localhost:3000/books.jsonld \
    --get \
    --header 'Preload: /hydra:member/*/author' \
    --header 'Fields: /hydra:member/*/author' \
    --header 'Fields: /hydra:member/*/title' \
    --verbose \
    --insecure 
```

## Protocol

The protocol is written in Markdown, compatible with [Mmark](https://github.com/mmarkdown/mmark).
It is then converted in the [the "xml2rfc" Version 3 Vocabulary](https://tools.ietf.org/html/rfc7991).

To contribute to the protocol itself:

* Make your changes
* [Download Mmark](https://github.com/mmarkdown/mmark/releases)
* [Download `xml2rfc` using pip](https://pypi.org/project/xml2rfc/): `pip install xml2rfc`
* Format the Markdown file: `mmark -markdown -w spec/vulcain.md`
* Generate the XML file: `mmark spec/vulcain.md > spec/vulcain.xml`
* Add the `docName` attribute to the `<rfc>` element (example: `docName="draft-dunglas-vulcain-01"`)
* Validate the generated XML file and generate the text file: `xml2rfc --text --v3 spec/vulcain.xml`
* Remove non-ASCII characters from the generated `vulcain.txt` file (example: K**é**vin)
* If appropriate, be sure to update the reference implementation accordingly
