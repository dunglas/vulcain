#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

assertRequests([
    ['/books.jsonld', ['headers' => ['Preload' => '"/hydra:member/*/author"']]],
    ['/books/1.jsonld', ['headers' => ['Preload' => '"/author"']]],
    ['/books/2.jsonld', ['headers' => ['Preload' => '"/author"']]],
    ['/authors/1.jsonld'],
], [
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/1.jsonld"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/2.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/authors/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books/2.jsonld"',
]);
