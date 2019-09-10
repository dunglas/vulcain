#!/usr/bin/env php
<?php

require __DIR__.'/vendor/autoload.php';

use Symfony\Component\HttpClient\CurlHttpClient;
use Symfony\Contracts\HttpClient\ResponseInterface;
use Psr\Log\AbstractLogger;

function printResponse(ResponseInterface $response): void
{
    echo "=======\n";
    echo "Headers\n";
    echo "=======\n";
    echo var_dump($response->getInfo());
    echo "=======\n";
    echo "Content\n";
    echo "=======\n";
    echo $response->getContent();
    echo "\n-------\n";
}

$gatewayUrl = $_SERVER['GATEWAY_URL'] ?? 'https://localhost:3000';

$logger = new class() extends AbstractLogger {
    public $logs = [];

    public function log($level, $message, array $context = [])
    {
        $this->logs[] = $message;
    }
};

$client = new CurlHttpClient(['verify_peer' => false, 'verify_host' => false, 'headers' => ['Cookie' => 'myCookie=bar']], 6, 5);
$client->setLogger($logger);

$urls = [
    '/books.jsonld?preload=/hydra:member/*/author',
    '/books/1.jsonld?preload=%2Fauthor',
    '/authors/1.jsonld',
];
foreach ($urls as $url) {
    $res = $client->request('GET', $gatewayUrl.$url);
    printResponse($res);
}

echo implode("\n", $logger->logs)."\n";

$expected = [
    'Request: "GET '.$gatewayUrl.'/books.jsonld?preload=/hydra:member/*/author"',
    'Queueing pushed response: "'.$gatewayUrl.'/books/1.jsonld?preload=%2Fauthor"',
    'Queueing pushed response: "'.$gatewayUrl.'/books/2.jsonld?preload=%2Fauthor"',
    'Queueing pushed response: "'.$gatewayUrl.'/authors/1.jsonld"',
    'Queueing pushed response: "'.$gatewayUrl.'/authors/1.jsonld"',
    'Response: "200 '.$gatewayUrl.'/books.jsonld?preload=/hydra:member/*/author"',
    'Accepting pushed response: "GET '.$gatewayUrl.'/books/1.jsonld?preload=%2Fauthor"',
    'Response: "200 '.$gatewayUrl.'/books/1.jsonld?preload=%2Fauthor"',
    'Accepting pushed response: "GET '.$gatewayUrl.'/authors/1.jsonld"',
    'Response: "200 '.$gatewayUrl.'/authors/1.jsonld"',
];

if ($logger->logs !== $expected) {
    fwrite(STDERR, "-".implode("\n-", array_diff($expected, $logger->logs))."\n");
    exit(1);
}
