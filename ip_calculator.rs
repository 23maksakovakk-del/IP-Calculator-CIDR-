// ip_calculator.rs — Rust версия

use std::env;
use std::fs;
use std::net::{Ipv4Addr, Ipv4Network};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
struct IPInfo {
    ip: String,
    mask: String,
    cidr: String,
    network: String,
    broadcast: String,
    host_min: String,
    host_max: String,
    num_hosts: u32,
    version: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Usage: cargo run -- <IP/CIDR>");
        println!("Пример: cargo run -- 192.168.1.0/24");
        std::process::exit(1);
    }

    let cidr = &args[1];
    let network: Ipv4Network = cidr.parse()?;

    let ip = network.network().to_string();
    let mask = network.mask().to_string();
    let cidr_str = format!("/{}", network.prefix());
    let network_addr = network.network().to_string();
    let broadcast = network.broadcast().to_string();

    let hosts: Vec<Ipv4Addr> = network.hosts().collect();
    let host_min = if hosts.is_empty() { "N/A".to_string() } else { hosts[0].to_string() };
    let host_max = if hosts.is_empty() { "N/A".to_string() } else { hosts[hosts.len()-1].to_string() };
    let num_hosts = if network.prefix() >= 31 { network.hosts().count() as u32 } else { network.hosts().count() as u32 };

    let info = IPInfo {
        ip,
        mask,
        cidr: cidr_str,
        network: network_addr,
        broadcast,
        host_min,
        host_max,
        num_hosts,
        version: "IPv4".to_string(),
    };

    println!("\x1b[36m🌐 IP Calculator (CIDR) (Rust)\x1b[0m");
    println!("Входные данные: {}", cidr);
    println!();
    println!("\x1b[32m─────────────────────────────────────────\x1b[0m");
    println!("IP-адрес:       {}", info.ip);
    println!("Маска (CIDR):   {} ({})", info.cidr, info.mask);
    println!("Сетевой адрес:  {}", info.network);
    println!("Широковещательный: {}", info.broadcast);
    println!("Диапазон хостов: {} — {}", info.host_min, info.host_max);
    println!("Количество хостов: {}", info.num_hosts);
    println!("Версия:         {}", info.version);
    println!("\x1b[32m─────────────────────────────────────────\x1b[0m");

    let json = serde_json::to_string_pretty(&info)?;
    fs::write("ip_calc_output.json", json)?;
    println!("\x1b[32m💾 Сохранено: ip_calc_output.json\x1b[0m");

    Ok(())
}
