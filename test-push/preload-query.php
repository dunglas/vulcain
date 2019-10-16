#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

assertRequests([
    '/books.jsonld?preload=/hydra:member/*/author',
    '/books/1.jsonld?preload=%2Fauthor',
    '/books/2.jsonld?preload=%2Fauthor',
    '/authors/1.jsonld',
], [
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/1.jsonld?preload=%2Fauthor"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/2.jsonld?preload=%2Fauthor"',
    'Queueing pushed response: "' . GATEWAY_URL . '/authors/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld?preload=%2Fauthor"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld?preload=%2Fauthor"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld?preload=/hydra:member/*/author"',
    'Response: "200 ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld?preload=/hydra:member/*/author"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld?preload=%2Fauthor"',
    'Response: "200 ' . GATEWAY_URL . '/books/2.jsonld?preload=%2Fauthor"',
]);
