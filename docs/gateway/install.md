# Install The Gateway Server

The Vulcain Gateway Server can be put in front of **any existing REST API** to transform it in a Vulcain-enabled API.
It works with hypermedia APIs ([JSON-LD](https://json-ld.org), [JSON:API](https://jsonapi.org/), [HAL](https://tools.ietf.org/html/draft-kelly-json-hal), [Siren](https://github.com/kevinswiber/siren) ...) as well as [with other non-hypermedia APIs by configuring the server using a subset of the OpenAPI specification](openapi.md).

Tip: the easiest way to create a hypermedia API is to use [the API Platform framework](https://api-platform.com) (by the same author than Vulcain).

**The Gateway Server is still experimental**

## Prebuilt Binary

First, download the archive corresponding to your operating system and architecture [from the release page](https://github.com/dunglas/vulcain/releases), extract the archive and open a shell in the resulting directory.

Note: Mac OS users must use the `Darwin` binary.

To use HTTP/2 Server Push, the connection must be encrypted with HTTPS.
To test the hub locally, use [OpenSSL](https://www.openssl.org/) ([Windows binaries](https://wiki.openssl.org/index.php/Binaries)) to generate a self-signed certificate:

    mkdir tls
    openssl req -x509 -newkey rsa:4096 -keyout tls/key.pem -out tls/cert.pem -days 365    

Then, on UNIX, run:

    UPSTREAM='http://your-api' ADDR=':3000' KEY_FILE='tls/key.pem' CERT_FILE='tls/cert.pem' ./vulcain

On Windows, start [PowerShell](https://docs.microsoft.com/en-us/powershell/), go into the extracted directory and run:

    $env:UPSTREAM='http://your-api'; $env:ADDR='localhost:3000'; $env:KEY_FILE='key.pem'; $env:CERT_FILE='cert.pem'; .\vulcain.exe

The Windows Defender Firewall will ask you if you want to allow `vulcain.exe` to communicate through it. Allow it for both public and private networks. If you use an antivirus, or another firewall software, be sure to whitelist `vulcain.exe`. 

The gateway is now available on `https://localhost:3000`.

To run it in production mode, and generate automatically a Let's Encrypt TLS certificate, just run the following command as root:

    UPSTREAM='http://your-api' ACME_HOSTS='example.com' ./vulcain

Using Windows in production is not recommended.

The value of the `ACME_HOSTS` environment variable must be updated to match your domain name(s).
A Let's Encrypt TLS certificate will be automatically generated.
If you omit this variable, the server will be exposed using a not encrypted HTTP connection, so you will not be able to use Server Push.

To compile the development version and register the demo page, see [CONTRIBUTING](../../CONTRIBUTING.md#start-a-demo-api-and-contribute-to-the-gateway-server).

## Docker Image

A Docker image is available on Docker Hub. The following command is enough to get a working server:

    docker run \
        -v /your/tls/certs:/tls \
        -e UPSTREAM='http://your-api' -e KEY_FILE='tls/key.pem' -e CERT_FILE='tls/cert.pem' \
        -p 443:443 \
        dunglas/vulcain

The gateway is available on `https://localhost`.

In production, run:

    docker run \
        -e UPSTREAM='http://your-api' -e ACME_HOSTS='example.com' \
        -p 80:80 -p 443:443 \
        dunglas/vulcain

Be sure to update the value of `ACME_HOSTS` to match your domain name(s), a Let's Encrypt TLS certificate will be automatically generated.

* [Configuration options](config.md)
* [Mapping a non-hypermedia API using OpenAPI](openapi.md)
