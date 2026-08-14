// ip_calculator.go — Go версия

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type IPInfo struct {
	IP         string `json:"ip"`
	Mask       string `json:"mask"`
	CIDR       string `json:"cidr"`
	Network    string `json:"network"`
	Broadcast  string `json:"broadcast"`
	HostMin    string `json:"host_min"`
	HostMax    string `json:"host_max"`
	NumHosts   uint64 `json:"num_hosts"`
	Version    string `json:"version"`
	MaskBinary string `json:"mask_binary"`
}

func calculate(cidr string) (*IPInfo, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ones, bits := ipnet.Mask.Size()
	mask := net.IP(ipnet.Mask).String()

	// Вычисляем широковещательный адрес
	broadcast := make(net.IP, len(ipnet.IP))
	copy(broadcast, ipnet.IP)
	for i := range broadcast {
		broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
	}

	// Вычисляем первый и последний хост
	first := make(net.IP, len(ipnet.IP))
	copy(first, ipnet.IP)
	last := make(net.IP, len(broadcast))
	copy(last, broadcast)

	if bits == 0 {
		// /0 — специальный случай
		return &IPInfo{
			IP:        ip.String(),
			Mask:      mask,
			CIDR:      fmt.Sprintf("/%d", ones),
			Network:   ipnet.IP.String(),
			Broadcast: "N/A",
			HostMin:   "N/A",
			HostMax:   "N/A",
			NumHosts:  0,
			Version:   "IPv4",
		}, nil
	}

	// Для подсетей с более чем 2 адресами
	if bits < 32 {
		first[3] = ipnet.IP[3] + 1
		last[3] = broadcast[3] - 1
	} else {
		// /32 — один хост
		first = ip
		last = ip
	}

	numHosts := uint64(1) << (uint(bits) - 1)
	if bits < 32 && bits > 0 {
		numHosts = (1 << (32 - uint(ones))) - 2
	} else if bits == 32 {
		numHosts = 1
	}

	return &IPInfo{
		IP:         ip.String(),
		Mask:       mask,
		CIDR:       fmt.Sprintf("/%d", ones),
		Network:    ipnet.IP.String(),
		Broadcast:  broadcast.String(),
		HostMin:    first.String(),
		HostMax:    last.String(),
		NumHosts:   numHosts,
		Version:    "IPv4",
		MaskBinary: maskToBinary(mask),
	}, nil
}

func maskToBinary(mask string) string {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return "N/A"
	}
	var result []string
	for _, p := range parts {
		i, _ := strconv.Atoi(p)
		result = append(result, fmt.Sprintf("%08b", i))
	}
	return strings.Join(result, ".")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ip_calculator.go <IP/CIDR>")
		fmt.Println("Пример: go run ip_calculator.go 192.168.1.0/24")
		os.Exit(1)
	}

	cidr := os.Args[1]
	info, err := calculate(cidr)
	if err != nil {
		fmt.Printf("\x1b[31m❌ Ошибка: %v\x1b[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("\x1b[36m🌐 IP Calculator (CIDR) (Go)\x1b[0m")
	fmt.Printf("Входные данные: %s\n", cidr)
	fmt.Println()
	fmt.Println("\x1b[32m─────────────────────────────────────────\x1b[0m")
	fmt.Printf("IP-адрес:       %s\n", info.IP)
	fmt.Printf("Маска (CIDR):   %s (%s)\n", info.CIDR, info.Mask)
	fmt.Printf("Сетевой адрес:  %s\n", info.Network)
	fmt.Printf("Широковещательный: %s\n", info.Broadcast)
	fmt.Printf("Диапазон хостов: %s — %s\n", info.HostMin, info.HostMax)
	fmt.Printf("Количество хостов: %d\n", info.NumHosts)
	fmt.Printf("Версия:         %s\n", info.Version)
	fmt.Printf("Маска (двоичная): %s\n", info.MaskBinary)
	fmt.Println("\x1b[32m─────────────────────────────────────────\x1b[0m")

	// Сохраняем JSON
	jsonData, _ := json.MarshalIndent(info, "", "  ")
	os.WriteFile("ip_calc_output.json", jsonData, 0644)
	fmt.Printf("\x1b[32m💾 Сохранено: ip_calc_output.json\x1b[0m\n")
}
