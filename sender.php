<?php
$host   = getenv('HOST') ?: '127.0.0.1';
$port = getenv('PORT') ?: '1961';

const baseBatchSize = 50;
const types = ['ALRM', 'FLOD', 'FIRE', 'QAKE', 'DOOM', 'TEST', 'PING'];
const levels = ['CRIT', 'EROR', 'WARN', 'INFO', 'TEST'];
const descrs = ['something happened', 'you hear noise', 'it\'s all different', 'ignore this message', 'silly propaganda'];

function randomEvent() {
    $type = types[array_rand(types)];
    $level = levels[array_rand(levels)];
    $time = date('Ymd\THis\Z000');
    $event = 'Event-' . strtolower("$level-$type");
    $eventid = rand(1, 31) . '-' . rand(0, 99);
    $descr = descrs[array_rand(descrs)];
    return "* * * $type $level $time\n$event $eventid\n$descr\nEND";
}

echo "Spamming at $host:$port\n";
echo "Ctrl-C to terminate\n";

$socket = socket_create(AF_INET, SOCK_DGRAM, SOL_UDP);

$count = 0;
$timeStart = time();

while (1) {
    $batchSize = rand(baseBatchSize, baseBatchSize * 2);
    for ($i = 0; $i < $batchSize; $i++) {
        $message = randomEvent();
        socket_sendto($socket, $message, strlen($message), 0, $host, $port);
    }
    $count += $batchSize;
    $time = time() - $timeStart;
    echo "Messages sent: $count, time from start: $time\n";
    sleep(1);
}

