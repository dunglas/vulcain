#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

use \Symfony\Contracts\HttpClient\ResponseInterface;

assertRequests([
    ['/books.jsonld?fields=/hydra:member/*/author&preload=/hydra:member/*/author', [], function (ResponseInterface $response) {
        $expectedContent = '{"hydra:member":["/books/1.jsonld?fields=%2Fauthor\u0026preload=%2Fauthor","/books/2.jsonld?fields=%2Fauthor\u0026preload=%2Fauthor"]}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
    ['/books/1.jsonld?fields=%2Fauthor', [], function (ResponseInterface $response) {
        $expectedContent = '{"author":"/authors/1.jsonld"}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
], [
    'Queueing pushed response: "' . GATEWAY_URL . '/authors/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld?fields=%2Fauthor&preload=%2Fauthor"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld?fields=%2Fauthor&preload=%2Fauthor"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld?fields=/hydra:member/*/author&preload=/hydra:member/*/author"',
    'Request: "GET ' . GATEWAY_URL . '/books/1.jsonld?fields=%2Fauthor"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld?fields=/hydra:member/*/author&preload=/hydra:member/*/author"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld?fields=%2Fauthor"',
]);
