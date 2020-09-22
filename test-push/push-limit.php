#!/usr/bin/env php
<?php

/**
 * The Push Limit must be configured to 2 before running this test.
 * See server_test.go::TestH2PushLimit
 */

require __DIR__ . '/tester.php';

assertRequests([
    ['/books.jsonld', ['headers' => ['Preload' => '"/hydra:member/*/author"']]],
], function (string $logs): int {
    $nb = preg_match_all('/Queueing pushed response/', $logs);
    
    return $nb == 2 ? 0 : $nb;
});
