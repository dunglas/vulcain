#!/usr/bin/env php
<?php

require __DIR__.'/vendor/autoload.php';

$gatewayUrl = $_SERVER['GATEWAY_URL'] ?? 'https://localhost:3000';

use Symfony\Component\HttpClient\CurlHttpClient;
use Psr\Log\AbstractLogger;

$logger = new class() extends AbstractLogger {
    public $logs = [];

    public function log($level, $message, array $context = [])
    {
        $this->logs[] = $message;
    }
};

$client = new CurlHttpClient(['verify_peer' => false, 'verify_host' => false, 'headers' => ['Cookie' => 'myCookie=bar']], 6, 5);
$client->setLogger($logger);

$books = $client->request('GET', "$gatewayUrl/books.jsonld?preload=/hydra:member/*/author");
$books->getHeaders();

$book1 = $client->request('GET', "$gatewayUrl/books/1.jsonld");
$book1->getHeaders();

$book2 = $client->request('GET', "$gatewayUrl/books/2.jsonld");
$book2->getHeaders();

$authors1 = $client->request('GET', "$gatewayUrl/authors/1.jsonld");
$authors1->getHeaders();

$authors2 = $client->request('GET', "$gatewayUrl/authors/1.jsonld");
$authors2->getHeaders();

fwrite(STDERR, implode("\n", $logger->logs)."\n");

$expected = [
    "Request: \"GET $gatewayUrl/books.jsonld?preload=/hydra:member/*/author\"",
    "Queueing pushed response: \"$gatewayUrl/books/1.jsonld?preload=%2Fauthor\"",
    "Queueing pushed response: \"$gatewayUrl/books/2.jsonld?preload=%2Fauthor\"",
    "Queueing pushed response: \"$gatewayUrl/authors/1.jsonld\"",
    "Queueing pushed response: \"$gatewayUrl/authors/1.jsonld\"",
    "Response: \"200 $gatewayUrl/books.jsonld?preload=/hydra:member/*/author\"",
    "Request: \"GET $gatewayUrl/books/1.jsonld\"",
    "Response: \"200 $gatewayUrl/books/1.jsonld\"",
    "Request: \"GET $gatewayUrl/books/2.jsonld\"",
    "Response: \"200 $gatewayUrl/books/2.jsonld\"",
    "Accepting pushed response: \"GET $gatewayUrl/authors/1.jsonld\"",
    "Response: \"200 $gatewayUrl/authors/1.jsonld\"",
    "Request: \"GET $gatewayUrl/authors/1.jsonld\"",
    "Response: \"200 $gatewayUrl/authors/1.jsonld\"",
];

if ($logger->logs !== $expected) {
    fwrite(STDERR, "-".implode("\n-", array_diff($expected, $logger->logs))."\n");
    exit(1);
}
