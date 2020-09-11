#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

use \Symfony\Contracts\HttpClient\ResponseInterface;

assertRequests([
    ['/books.jsonld?fields="/hydra:member/*/author"&preload="/hydra:member/*/author"', [], function (ResponseInterface $response) {
        $expectedContent = '{"hydra:member":["/books/1.jsonld?fields=%22%2Fauthor%22\u0026preload=%22%2Fauthor%22","/books/2.jsonld?fields=%22%2Fauthor%22\u0026preload=%22%2Fauthor%22"]}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
    ['/books/1.jsonld?fields="%2Fauthor"', [], function (ResponseInterface $response) {
        $expectedContent = '{"author":"/authors/1.jsonld"}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
], [
    'Queueing pushed response: "' . GATEWAY_URL . '/authors/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld?fields=%22%2Fauthor%22&preload=%22%2Fauthor%22"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld?fields=%22%2Fauthor%22&preload=%22%2Fauthor%22"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld?fields=%22/hydra:member/*/author%22&preload=%22/hydra:member/*/author%22"',
    'Request: "GET ' . GATEWAY_URL . '/books/1.jsonld?fields=%22%2Fauthor%22"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld?fields=%22/hydra:member/*/author%22&preload=%22/hydra:member/*/author%22"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld?fields=%22%2Fauthor%22"',
]);
