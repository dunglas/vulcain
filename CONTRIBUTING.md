# Contributing

## License and Copyright Attribution

When you open a Pull Request to the project, you agree to license your code under the [GNU AFFERO GENERAL PUBLIC LICENSE](LICENSE)
and to transfer the copyright on the submitted code to [Kévin Dunglas](https://dunglas.fr).

Be sure to you have the right to do that (if you are a professional, ask your company)!

If you include code from another project, please mention it in the Pull Request description and credit the original author.

## Start a Demo API and Contribute to the Gateway Server

Clone the project:

    $ git clone https://github.com/dunglas/vulcain
    $ cd vulcain

Install the dependencies:

    $ go get

Run the server:

    $ go run cmd/vulcain/main.go

Run the fixture API:

    # You must run the server too
    $ cd fixtures/
    $ go run main.go

Go to `https://localhost:3000` and accept the self-signed certificate.
Go on `http://localhost:8081` and enjoy!

An API using an OpenAPI mapping is available on `https://localhost:3000/oa/books.json`.

To run the test suite:

    $ go test -v -timeout 30s github.com/dunglas/vulcain/gateway

### curl Examples

Preload all relations referenced in the `hydra:member`, then in the author relationship, but only include the title and the author of these relations:

```
curl https://localhost:3000/books.jsonld \
    --get \
    --data 'preload="/hydra:member/*/author"' \
    --data 'fields="/hydra:member/*/author", "/hydra:member/*/title"' \
    --verbose \
    --insecure
```

Using headers:

```
curl https://localhost:3000/books.jsonld \
    --get \
    --header 'Preload: "/hydra:member/*/author"' \
    --header 'Fields: "/hydra:member/*/author", "/hydra:member/*/title"' \
    --verbose \
    --insecure
```

## Protocol

The protocol is written in Markdown, compatible with [Mmark](https://mmark.miek.nl/).
It is then converted in the [the "xml2rfc" Version 3 Vocabulary](https://tools.ietf.org/html/rfc7991).

To contribute to the protocol itself:

* Make your changes
* [Download Mmark](https://github.com/mmarkdown/mmark/releases)
* [Download `xml2rfc` using pip](https://pypi.org/project/xml2rfc/): `pip install xml2rfc`
* Format the Markdown file: `mmark -markdown -w spec/vulcain.md`
* Generate the XML file: `mmark spec/vulcain.md > spec/vulcain.xml`
* Validate the generated XML file and generate the text file: `xml2rfc --text --v3 spec/vulcain.xml`
* Remove non-ASCII characters from the generated `vulcain.txt` file (example: K**é**vin, Andr**é**, **Ã**elik)
* If appropriate, be sure to update the reference implementation accordingly
