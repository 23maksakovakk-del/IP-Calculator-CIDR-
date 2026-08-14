<?php
// ip_calculator.php — PHP версия

class IPCalculator {
    private $cidr;
    private $ip;
    private $prefix;
    private $mask;
    private $network;
    private $broadcast;
    private $hostMin;
    private $hostMax;
    private $numHosts;

    public function __construct($cidr) {
        $this->cidr = $cidr;
        $parts = explode('/', $cidr);
        $this->ip = $parts[0];
        $this->prefix = (int)$parts[1];
        $this->mask = $this->calculateMask($this->prefix);
        $this->network = $this->calculateNetwork();
        $this->broadcast = $this->calculateBroadcast();
        $this->hostMin = $this->calculateHostMin();
        $this->hostMax = $this->calculateHostMax();
        $this->numHosts = pow(2, 32 - $this->prefix) - 2;
        if ($this->prefix == 32) $this->numHosts = 1;
        if ($this->prefix == 0) $this->numHosts = 0;
    }

    private function calculateMask($prefix) {
        $maskInt = 0xFFFFFFFF << (32 - $prefix);
        return sprintf("%d.%d.%d.%d",
            ($maskInt >> 24) & 0xFF,
            ($maskInt >> 16) & 0xFF,
            ($maskInt >> 8) & 0xFF,
            $maskInt & 0xFF);
    }

    private function ipToInt($ip) {
        $parts = explode('.', $ip);
        return ($parts[0] << 24) | ($parts[1] << 16) | ($parts[2] << 8) | $parts[3];
    }

    private function intToIp($int) {
        return sprintf("%d.%d.%d.%d",
            ($int >> 24) & 0xFF,
            ($int >> 16) & 0xFF,
            ($int >> 8) & 0xFF,
            $int & 0xFF);
    }

    private function calculateNetwork() {
        $ipInt = $this->ipToInt($this->ip);
        $maskInt = $this->ipToInt($this->mask);
        return $this->intToIp($ipInt & $maskInt);
    }

    private function calculateBroadcast() {
        $networkInt = $this->ipToInt($this->network);
        $maskInt = $this->ipToInt($this->mask);
        return $this->intToIp($networkInt | (~$maskInt & 0xFFFFFFFF));
    }

    private function calculateHostMin() {
        if ($this->prefix == 32) return $this->ip;
        if ($this->prefix == 0) return "N/A";
        $networkInt = $this->ipToInt($this->network);
        return $this->intToIp($networkInt + 1);
    }

    private function calculateHostMax() {
        if ($this->prefix == 32) return $this->ip;
        if ($this->prefix == 0) return "N/A";
        $broadcastInt = $this->ipToInt($this->broadcast);
        return $this->intToIp($broadcastInt - 1);
    }

    public function printInfo() {
        echo "\033[36m🌐 IP Calculator (CIDR) (PHP)\033[0m\n";
        echo "Входные данные: {$this->cidr}\n";
        echo "\n";
        echo "\033[32m─────────────────────────────────────────\033[0m\n";
        echo "IP-адрес:       {$this->ip}\n";
        echo "Маска (CIDR):   /{$this->prefix} ({$this->mask})\n";
        echo "Сетевой адрес:  {$this->network}\n";
        echo "Широковещательный: {$this->broadcast}\n";
        echo "Диапазон хостов: {$this->hostMin} — {$this->hostMax}\n";
        echo "Количество хостов: {$this->numHosts}\n";
        echo "Версия:         IPv4\n";
        echo "\033[32m─────────────────────────────────────────\033[0m\n";
    }

    public function saveJSON($filename = 'ip_calc_output.json') {
        $data = [
            'ip' => $this->ip,
            'mask' => $this->mask,
            'cidr' => "/{$this->prefix}",
            'network' => $this->network,
            'broadcast' => $this->broadcast,
            'host_min' => $this->hostMin,
            'host_max' => $this->hostMax,
            'num_hosts' => $this->numHosts,
            'version' => 'IPv4'
        ];
        file_put_contents($filename, json_encode($data, JSON_PRETTY_PRINT));
        echo "\033[32m💾 Сохранено: {$filename}\033[0m\n";
    }
}

function main($argv) {
    if ($argc < 2) {
        echo "Usage: php ip_calculator.php <IP/CIDR>\n";
        echo "Пример: php ip_calculator.php 192.168.1.0/24\n";
        exit(1);
    }

    $cidr = $argv[1];
    try {
        $calc = new IPCalculator($cidr);
        $calc->printInfo();
        $calc->saveJSON();
    } catch (Exception $e) {
        echo "\033[31m❌ Ошибка: " . $e->getMessage() . "\033[0m\n";
        exit(1);
    }
}

$argc = $_SERVER['argc'] ?? 0;
$argv = $_SERVER['argv'] ?? [];
main($argv);
?>
