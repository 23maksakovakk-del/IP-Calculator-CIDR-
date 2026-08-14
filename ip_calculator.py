

### 1. `ip_calculator.py` (Python)

```python
# ip_calculator.py — Python версия

import sys
import json
import ipaddress
from colorama import init, Fore, Style

init(autoreset=True)

class IPCalculator:
    def __init__(self, cidr):
        self.cidr = cidr
        self.network = ipaddress.ip_network(cidr, strict=False)
        self.ip = self.network.network_address
        self.mask = self.network.netmask
        self.broadcast = self.network.broadcast_address
        self.hosts = list(self.network.hosts())
        self.num_hosts = self.network.num_addresses - 2 if self.network.num_addresses > 2 else 0
        self.version = self.network.version

    def get_info(self):
        info = {
            "ip": str(self.ip),
            "mask": str(self.mask),
            "cidr": f"/{self.network.prefixlen}",
            "network": str(self.network.network_address),
            "broadcast": str(self.broadcast) if self.broadcast else "N/A",
            "host_min": str(self.hosts[0]) if self.hosts else "N/A",
            "host_max": str(self.hosts[-1]) if self.hosts else "N/A",
            "num_hosts": self.num_hosts,
            "version": f"IPv{self.version}",
            "netmask_binary": self._mask_to_binary(self.mask),
        }
        return info

    def _mask_to_binary(self, mask):
        if self.version == 4:
            return '.'.join([bin(int(o))[2:].zfill(8) for o in str(mask).split('.')])
        else:
            return "N/A (IPv6)"

    def print_info(self):
        info = self.get_info()
        print(Fore.CYAN + "🌐 IP Calculator (CIDR) (Python)")
        print(f"Входные данные: {self.cidr}")
        print()
        print(Fore.GREEN + "─────────────────────────────────────────")
        print(f"IP-адрес:       {info['ip']}")
        print(f"Маска (CIDR):   {info['cidr']} ({info['mask']})")
        print(f"Сетевой адрес:  {info['network']}")
        print(f"Широковещательный: {info['broadcast']}")
        print(f"Диапазон хостов: {info['host_min']} — {info['host_max']}")
        print(f"Количество хостов: {info['num_hosts']}")
        print(f"Версия:         {info['version']}")
        if info['netmask_binary'] != "N/A (IPv6)":
            print(f"Маска (двоичная): {info['netmask_binary']}")
        print(Fore.GREEN + "─────────────────────────────────────────")

    def save_json(self, filename="ip_calc_output.json"):
        info = self.get_info()
        with open(filename, 'w') as f:
            json.dump(info, f, indent=2)
        print(Fore.GREEN + f"💾 Сохранено: {filename}")

def main():
    if len(sys.argv) < 2:
        print("Usage: python ip_calculator.py <IP/CIDR>")
        print("Пример: python ip_calculator.py 192.168.1.0/24")
        sys.exit(1)

    cidr = sys.argv[1]
    try:
        calc = IPCalculator(cidr)
        calc.print_info()
        calc.save_json()
    except ValueError as e:
        print(Fore.RED + f"❌ Ошибка: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
