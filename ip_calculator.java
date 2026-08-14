// ip_calculator.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.net.*;
import java.util.*;

public class ip_calculator {
    private String cidr;
    private String ip;
    private int prefix;
    private String mask;
    private String network;
    private String broadcast;
    private String hostMin;
    private String hostMax;
    private long numHosts;

    public ip_calculator(String cidr) throws Exception {
        this.cidr = cidr;
        String[] parts = cidr.split("/");
        this.ip = parts[0];
        this.prefix = Integer.parseInt(parts[1]);
        this.mask = calculateMask(this.prefix);
        this.network = calculateNetwork();
        this.broadcast = calculateBroadcast();
        this.hostMin = calculateHostMin();
        this.hostMax = calculateHostMax();
        this.numHosts = (long)Math.pow(2, 32 - this.prefix) - 2;
        if (this.prefix == 32) this.numHosts = 1;
        if (this.prefix == 0) this.numHosts = 0;
    }

    private String calculateMask(int prefix) {
        long maskInt = 0xFFFFFFFFL << (32 - prefix);
        return String.format("%d.%d.%d.%d",
            (maskInt >> 24) & 0xFF,
            (maskInt >> 16) & 0xFF,
            (maskInt >> 8) & 0xFF,
            maskInt & 0xFF);
    }

    private long ipToLong(String ip) throws Exception {
        String[] parts = ip.split("\\.");
        long result = 0;
        for (int i = 0; i < 4; i++) {
            result = (result << 8) | Long.parseLong(parts[i]);
        }
        return result;
    }

    private String longToIp(long ip) {
        return String.format("%d.%d.%d.%d",
            (ip >> 24) & 0xFF,
            (ip >> 16) & 0xFF,
            (ip >> 8) & 0xFF,
            ip & 0xFF);
    }

    private String calculateNetwork() throws Exception {
        long ipLong = ipToLong(ip);
        long maskLong = ipToLong(mask);
        return longToIp(ipLong & maskLong);
    }

    private String calculateBroadcast() throws Exception {
        long networkLong = ipToLong(network);
        long maskLong = ipToLong(mask);
        return longToIp(networkLong | (~maskLong & 0xFFFFFFFFL));
    }

    private String calculateHostMin() throws Exception {
        if (prefix == 32) return ip;
        if (prefix == 0) return "N/A";
        long networkLong = ipToLong(network);
        return longToIp(networkLong + 1);
    }

    private String calculateHostMax() throws Exception {
        if (prefix == 32) return ip;
        if (prefix == 0) return "N/A";
        long broadcastLong = ipToLong(broadcast);
        return longToIp(broadcastLong - 1);
    }

    public void printInfo() throws Exception {
        System.out.println("\u001B[36m🌐 IP Calculator (CIDR) (Java)\u001B[0m");
        System.out.println("Входные данные: " + cidr);
        System.out.println();
        System.out.println("\u001B[32m─────────────────────────────────────────\u001B[0m");
        System.out.println("IP-адрес:       " + ip);
        System.out.println("Маска (CIDR):   /" + prefix + " (" + mask + ")");
        System.out.println("Сетевой адрес:  " + network);
        System.out.println("Широковещательный: " + broadcast);
        System.out.println("Диапазон хостов: " + hostMin + " — " + hostMax);
        System.out.println("Количество хостов: " + numHosts);
        System.out.println("Версия:         IPv4");
        System.out.println("\u001B[32m─────────────────────────────────────────\u001B[0m");
    }

    public void saveJSON(String filename) throws IOException {
        Map<String, Object> info = new LinkedHashMap<>();
        info.put("ip", ip);
        info.put("mask", mask);
        info.put("cidr", "/" + prefix);
        info.put("network", network);
        info.put("broadcast", broadcast);
        info.put("host_min", hostMin);
        info.put("host_max", hostMax);
        info.put("num_hosts", numHosts);
        info.put("version", "IPv4");

        String json = new com.google.gson.GsonBuilder().setPrettyPrinting().create().toJson(info);
        Files.write(Paths.get(filename), json.getBytes());
        System.out.println("\u001B[32m💾 Сохранено: " + filename + "\u001B[0m");
    }

    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.out.println("Usage: java ip_calculator <IP/CIDR>");
            System.out.println("Пример: java ip_calculator 192.168.1.0/24");
            System.exit(1);
        }

        String cidr = args[0];
        try {
            ip_calculator calc = new ip_calculator(cidr);
            calc.printInfo();
            calc.saveJSON("ip_calc_output.json");
        } catch (Exception e) {
            System.out.println("\u001B[31m❌ Ошибка: " + e.getMessage() + "\u001B[0m");
            System.exit(1);
        }
    }
}
