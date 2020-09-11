#!/usr/bin/env php
<?php

/**
 * The Push Limit must be configured to 2 before running this test.
 * See server_test.go::TestH2PushLimit
 */

require __DIR__ . '/tester.php';

assertRequests([
    ['/books.jsonld', ['headers' => ['Preload' => '"/hydra:member/*/author"']]],
], [
    'Queueing pushed response: "' . GATEWAY_URL . '/books/1.jsonld"',
    'Queueing pushed response: "' . GATEWAY_URL . '/books/2.jsonld"',
    'Request: "GET ' . GATEWAY_URL . '/books.jsonld"',
    'Response: "200 ' . GATEWAY_URL . '/books.jsonld"',
]);
