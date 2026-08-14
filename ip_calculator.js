// ip_calculator.js — JavaScript версия

const fs = require('fs');

class IPCalculator {
    constructor(cidr) {
        this.cidr = cidr;
        const parts = cidr.split('/');
        this.ip = parts[0];
        this.prefix = parseInt(parts[1]);
        this.mask = this._calculateMask(this.prefix);
        this.network = this._calculateNetwork();
        this.broadcast = this._calculateBroadcast();
        this.hostMin = this._calculateHostMin();
        this.hostMax = this._calculateHostMax();
        this.numHosts = Math.pow(2, 32 - this.prefix) - 2;
        if (this.prefix === 32) this.numHosts = 1;
        if (this.prefix === 0) this.numHosts = 0;
    }

    _calculateMask(prefix) {
        const mask = [];
        let bits = prefix;
        for (let i = 0; i < 4; i++) {
            if (bits >= 8) {
                mask.push(255);
                bits -= 8;
            } else if (bits > 0) {
                mask.push(256 - Math.pow(2, 8 - bits));
                bits = 0;
            } else {
                mask.push(0);
            }
        }
        return mask.join('.');
    }

    _ipToInt(ip) {
        return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet), 0);
    }

    _intToIp(int) {
        return [(int >> 24) & 255, (int >> 16) & 255, (int >> 8) & 255, int & 255].join('.');
    }

    _calculateNetwork() {
        const ipInt = this._ipToInt(this.ip);
        const maskInt = this._ipToInt(this.mask);
        return this._intToIp(ipInt & maskInt);
    }

    _calculateBroadcast() {
        const networkInt = this._ipToInt(this.network);
        const maskInt = this._ipToInt(this.mask);
        return this._intToIp(networkInt | (~maskInt >>> 0));
    }

    _calculateHostMin() {
        if (this.prefix === 32) return this.ip;
        if (this.prefix === 0) return 'N/A';
        return this._intToIp(this._ipToInt(this.network) + 1);
    }

    _calculateHostMax() {
        if (this.prefix === 32) return this.ip;
        if (this.prefix === 0) return 'N/A';
        return this._intToIp(this._ipToInt(this.broadcast) - 1);
    }

    getInfo() {
        return {
            ip: this.ip,
            mask: this.mask,
            cidr: `/${this.prefix}`,
            network: this.network,
            broadcast: this.broadcast,
            host_min: this.hostMin,
            host_max: this.hostMax,
            num_hosts: this.numHosts,
            version: 'IPv4'
        };
    }

    printInfo() {
        const info = this.getInfo();
        console.log('\x1b[36m🌐 IP Calculator (CIDR) (JavaScript)\x1b[0m');
        console.log(`Входные данные: ${this.cidr}`);
        console.log();
        console.log('\x1b[32m─────────────────────────────────────────\x1b[0m');
        console.log(`IP-адрес:       ${info.ip}`);
        console.log(`Маска (CIDR):   ${info.cidr} (${info.mask})`);
        console.log(`Сетевой адрес:  ${info.network}`);
        console.log(`Широковещательный: ${info.broadcast}`);
        console.log(`Диапазон хостов: ${info.host_min} — ${info.host_max}`);
        console.log(`Количество хостов: ${info.num_hosts}`);
        console.log(`Версия:         ${info.version}`);
        console.log('\x1b[32m─────────────────────────────────────────\x1b[0m');
    }

    saveJSON(filename = 'ip_calc_output.json') {
        const info = this.getInfo();
        fs.writeFileSync(filename, JSON.stringify(info, null, 2));
        console.log(`\x1b[32m💾 Сохранено: ${filename}\x1b[0m`);
    }
}

function main() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.log('Usage: node ip_calculator.js <IP/CIDR>');
        console.log('Пример: node ip_calculator.js 192.168.1.0/24');
        process.exit(1);
    }

    const cidr = args[0];
    try {
        const calc = new IPCalculator(cidr);
        calc.printInfo();
        calc.saveJSON();
    } catch (err) {
        console.error(`\x1b[31m❌ Ошибка: ${err.message}\x1b[0m`);
        process.exit(1);
    }
}

if (require.main === module) main();
