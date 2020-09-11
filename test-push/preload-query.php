#!/usr/bin/env php
<?php

require __DIR__ . '/tester.php';

assertRequests([
    '/books.jsonld?preload="/hydra:member/*/author"',
    '/books/1.jsonld?preload="%2Fauthor"',
    '/books/2.jsonld?preload="%2Fauthor"',
    '/authors/1.jsonld',
], [
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/1.jsonld?preload=%22%2Fauthor%22"',
    'Accepting pushed response: "GET ' . GATEWAY_URL . '/books/2.jsonld?preload=%22%2Fauthor%22"',
    'Queueing pushed response: "' . GATEWAY_URL . '/authors/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld?preload=%22%2Fauthor%22"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld?preload=%22%2Fauthor%22"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld?preload=%22/hydra:member/*/author%22"',
    'Response: "200 ' . GATEWAY_URL . '/authors/1.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld?preload=%22/hydra:member/*/author%22"',
    'Response: "200 ' . GATEWAY_URL . '/books/1.jsonld?preload=%22%2Fauthor%22"',
    'Response: "200 ' . GATEWAY_URL . '/books/2.jsonld?preload=%22%2Fauthor%22"',
]);
