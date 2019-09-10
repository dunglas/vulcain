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

$client = new CurlHttpClient(['verify_peer' => false, 'verify_host' => false, 'headers' => ['Cookie' => 'myCookie=bar']], 6, 2);
$client->setLogger($logger);

$books = $client->request('GET', "$gatewayUrl/books.jsonld?preload=/hydra:member/*");
$books->getHeaders();

$book1 = $client->request('GET', "$gatewayUrl/books/1.jsonld");
$book1->getHeaders();

$book2 = $client->request('GET', "$gatewayUrl/books/2.jsonld");
$book2->getHeaders();

$expected = [
    "Request: \"GET $gatewayUrl/books.jsonld?preload=/hydra:member/*\"",
    "Queueing pushed response: \"$gatewayUrl/books/1.jsonld\"",
    "Queueing pushed response: \"$gatewayUrl/books/2.jsonld\"",
    "Response: \"200 $gatewayUrl/books.jsonld?preload=/hydra:member/*\"",
    "Accepting pushed response: \"GET $gatewayUrl/books/1.jsonld\"",
    "Response: \"200 $gatewayUrl/books/1.jsonld\"",
    "Accepting pushed response: \"GET $gatewayUrl/books/2.jsonld\"",
    "Response: \"200 $gatewayUrl/books/2.jsonld\"",
];

var_dump($logger->logs);

exit($logger->logs === $expected ? 0 : 1);
