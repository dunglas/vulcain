<?php

require __DIR__ . '/vendor/autoload.php';

use Psr\Log\AbstractLogger;
use SebastianBergmann\Diff\Differ;
use Symfony\Component\HttpClient\CurlHttpClient;
use Symfony\Contracts\HttpClient\ResponseInterface;

define('GATEWAY_URL', $_SERVER['GATEWAY_URL'] ?? 'https://localhost:3000');

function printResponse(ResponseInterface $response): void
{
    echo "=======\n";
    echo "Headers\n";
    echo "=======\n";
    var_dump($response->getInfo());
    echo "=======\n";
    echo "Content\n";
    echo "=======\n";
    echo $response->getContent();
    echo "\n-------\n";
}

/**
 * @param string[]|callable $expected
 */
function assertRequests(array $requests, array|callable $expected)
{
    $logger = new class() extends AbstractLogger {
        public $logs = [];

        public function log($level, \Stringable|string $message, array $context = []): void
        {
            $this->logs[] = $message;
        }
    };

    $client = new CurlHttpClient(['verify_peer' => false, 'verify_host' => false, 'headers' => ['Cookie' => 'myCookie=bar']], 6, 5);
    $client->setLogger($logger);

    foreach ($requests as $request) {
        $request = (array) $request;
        $res = $client->request('GET', GATEWAY_URL . $request[0], $request[1] ?? []);
        printResponse($res);

        if (isset($request[2])) {
            $request[2]($res);
        }
    }

    $logs = implode("\n", $logger->logs);
    echo $logs . "\n";

    if (is_array($expected)) {
        sort($logger->logs);
        if ($logger->logs !== $expected) {
            fwrite(STDERR, (new Differ())->diff(implode("\n", $expected), implode("\n", $logger->logs)));
            exit(1);
        }

        exit(0);
    }

    exit($expected($logs));
}
