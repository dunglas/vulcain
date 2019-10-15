#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

use \Symfony\Contracts\HttpClient\ResponseInterface;

assertRequests([
    ['/books.jsonld', ['headers' => ['Fields' => '/hydra:member/*/author', 'Vulcain-Explicit-Request' => 'attack']], function (ResponseInterface $response) {
        $expectedContent = '{"hydra:member":["/books/1.jsonld","/books/2.jsonld"]}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
    ['/books/1.jsonld', ['headers' => ['Fields' => '/author']], function (ResponseInterface $response) {
        $expectedContent = '{"author":"/authors/1.jsonld"}';
        $content = $response->getContent();
        if ($expectedContent !== $content) {
            throw new \UnexpectedValueException(sprintf('Expected "%s", got "%s"', $expectedContent, $content));
        }
    }],
], [
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld"',
    'Request: "GET ' . GATEWAY_URL . '/books/1.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld"',
]);
