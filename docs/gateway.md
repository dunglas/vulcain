# The Gateway Server

The Vulcain Gateway Server can be put in front of **any existing REST API** to transform it in a Vulcain-enabled API. It works with hypermedia APIs (JSON-LD, JSON:API, HAL...) as well as with other API by configuring the server using a subset of the OpenAPI specification.

## Prebuilt Binary

First, download the archive corresponding to your operating system and architecture [from the release page](https://github.com/dunglas/dunglas/releases), extract the archive and open a shell in the resulting directory.

Note: Mac OS users must use the `Darwin` binary.

Then, on UNIX, run:

    UPSTREAM='http://your-api' ./vulcain

On Windows, start [PowerShell](https://docs.microsoft.com/en-us/powershell/), go into the extracted directory and run:

    $env:UPSTREAM='http://your-api'; $env:ADDR='localhost:3000'; $env:DEMO='1'; $env:ALLOW_ANONYMOUS='1'; $env:CORS_ALLOWED_ORIGINS='*'; $env:PUBLISH_ALLOWED_ORIGINS='http://localhost:3000'; .\mercure.exe

The Windows Defender Firewall will ask you if you want to allow `mercure.exe` to communicate through it. Allow it for both public and private networks. If you use an antivirus, or another firewall software, be sure to whitelist `mercure.exe`. 

The server is now available on `http://localhost:3000`, with the demo mode enabled. Because `ALLOW_ANONYMOUS` is set to `1`, anonymous subscribers are allowed.

To run it in production mode, and generate automatically a Let's Encrypt TLS certificate, just run the following command as root:

    JWT_KEY='!ChangeMe!' ACME_HOSTS='example.com' ./mercure

Using Windows in production is not recommended.

The value of the `ACME_HOSTS` environment variable must be updated to match your domain name(s).
A Let's Encrypt TLS certificate will be automatically generated.
If you omit this variable, the server will be exposed using a not encrypted HTTP connection.

When the server is up and running, the following endpoints are available:

* `POST https://example.com/hub`: to publish updates
* `GET https://example.com/hub`: to subscribe to updates

See [the protocol](spec/mercure.md) for further informations.

To compile the development version and register the demo page, see [CONTRIBUTING.md](CONTRIBUTING.md#hub).

## Docker Image

A Docker image is available on Docker Hub. The following command is enough to get a working server in demo mode:

    docker run \
        -e JWT_KEY='!ChangeMe!' -e DEMO=1 -e ALLOW_ANONYMOUS=1 -e CORS_ALLOWED_ORIGINS=* -e PUBLISH_ALLOWED_ORIGINS='http://localhost' \
        -p 80:80 \
        dunglas/mercure

The server, in demo mode, is available on `http://localhost`. Anonymous subscribers are allowed.

In production, run:

    docker run \
        -e JWT_KEY='!ChangeMe!' -e ACME_HOSTS='example.com' \
        -p 80:80 -p 443:443 \
        dunglas/mercure

Be sure to update the value of `ACME_HOSTS` to match your domain name(s), a Let's Encrypt TLS certificate will be automatically generated.
