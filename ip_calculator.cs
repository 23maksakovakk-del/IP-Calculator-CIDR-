// ip_calculator.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Text.Json;

class IPCalculator {
    private string cidr;
    private string ip;
    private int prefix;
    private string mask;
    private string network;
    private string broadcast;
    private string hostMin;
    private string hostMax;
    private long numHosts;

    public IPCalculator(string cidr) {
        this.cidr = cidr;
        var parts = cidr.Split('/');
        this.ip = parts[0];
        this.prefix = int.Parse(parts[1]);
        this.mask = CalculateMask(this.prefix);
        this.network = CalculateNetwork();
        this.broadcast = CalculateBroadcast();
        this.hostMin = CalculateHostMin();
        this.hostMax = CalculateHostMax();
        this.numHosts = (long)Math.Pow(2, 32 - this.prefix) - 2;
        if (this.prefix == 32) this.numHosts = 1;
        if (this.prefix == 0) this.numHosts = 0;
    }

    private string CalculateMask(int prefix) {
        uint maskInt = 0xFFFFFFFFu << (32 - prefix);
        return $"{((maskInt >> 24) & 0xFF)}.{((maskInt >> 16) & 0xFF)}.{((maskInt >> 8) & 0xFF)}.{maskInt & 0xFF}";
    }

    private uint IpToUInt(string ip) {
        var parts = ip.Split('.');
        uint result = 0;
        for (int i = 0; i < 4; i++) {
            result = (result << 8) | uint.Parse(parts[i]);
        }
        return result;
    }

    private string UIntToIp(uint ip) {
        return $"{((ip >> 24) & 0xFF)}.{((ip >> 16) & 0xFF)}.{((ip >> 8) & 0xFF)}.{ip & 0xFF}";
    }

    private string CalculateNetwork() {
        uint ipLong = IpToUInt(ip);
        uint maskLong = IpToUInt(mask);
        return UIntToIp(ipLong & maskLong);
    }

    private string CalculateBroadcast() {
        uint networkLong = IpToUInt(network);
        uint maskLong = IpToUInt(mask);
        return UIntToIp(networkLong | (~maskLong));
    }

    private string CalculateHostMin() {
        if (prefix == 32) return ip;
        if (prefix == 0) return "N/A";
        uint networkLong = IpToUInt(network);
        return UIntToIp(networkLong + 1);
    }

    private string CalculateHostMax() {
        if (prefix == 32) return ip;
        if (prefix == 0) return "N/A";
        uint broadcastLong = IpToUInt(broadcast);
        return UIntToIp(broadcastLong - 1);
    }

    public void PrintInfo() {
        Console.WriteLine("\u001B[36m🌐 IP Calculator (CIDR) (C#)\u001B[0m");
        Console.WriteLine($"Входные данные: {cidr}");
        Console.WriteLine();
        Console.WriteLine("\u001B[32m─────────────────────────────────────────\u001B[0m");
        Console.WriteLine($"IP-адрес:       {ip}");
        Console.WriteLine($"Маска (CIDR):   /{prefix} ({mask})");
        Console.WriteLine($"Сетевой адрес:  {network}");
        Console.WriteLine($"Широковещательный: {broadcast}");
        Console.WriteLine($"Диапазон хостов: {hostMin} — {hostMax}");
        Console.WriteLine($"Количество хостов: {numHosts}");
        Console.WriteLine("Версия:         IPv4");
        Console.WriteLine("\u001B[32m─────────────────────────────────────────\u001B[0m");
    }

    public void SaveJSON(string filename) {
        var info = new {
            ip = this.ip,
            mask = this.mask,
            cidr = "/" + this.prefix,
            network = this.network,
            broadcast = this.broadcast,
            host_min = this.hostMin,
            host_max = this.hostMax,
            num_hosts = this.numHosts,
            version = "IPv4"
        };
        string json = JsonSerializer.Serialize(info, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"\u001B[32m💾 Сохранено: {filename}\u001B[0m");
    }

    public static void Main(string[] args) {
        if (args.Length < 1) {
            Console.WriteLine("Usage: dotnet run <IP/CIDR>");
            Console.WriteLine("Пример: dotnet run 192.168.1.0/24");
            return;
        }

        string cidr = args[0];
        try {
            var calc = new IPCalculator(cidr);
            calc.PrintInfo();
            calc.SaveJSON("ip_calc_output.json");
        } catch (Exception e) {
            Console.WriteLine($"\u001B[31m❌ Ошибка: {e.Message}\u001B[0m");
        }
    }
}
