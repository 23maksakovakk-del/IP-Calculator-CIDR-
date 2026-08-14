# ip_calculator.rb — Ruby версия

require 'ipaddr'
require 'json'

class IPCalculator
  def initialize(cidr)
    @cidr = cidr
    @ipaddr = IPAddr.new(cidr)
    @ip = @ipaddr.to_s
    @prefix = @ipaddr.prefix
    @mask = @ipaddr.inspect.split('/')[0]
    @network = @ipaddr.to_s
    @broadcast = @ipaddr.to_range.last.to_s
    @host_min = @ipaddr.to_range.first.to_s
    @host_max = @ipaddr.to_range.last.to_s
    @num_hosts = @ipaddr.to_range.size - 2
    if @prefix == 32
      @num_hosts = 1
      @broadcast = @ip
    elsif @prefix == 0
      @num_hosts = 0
      @broadcast = "N/A"
    end
  end

  def info
    {
      ip: @ip,
      mask: @mask,
      cidr: "/#{@prefix}",
      network: @network,
      broadcast: @broadcast,
      host_min: @host_min,
      host_max: @host_max,
      num_hosts: @num_hosts,
      version: "IPv4"
    }
  end

  def print_info
    info = self.info
    puts "\e[36m🌐 IP Calculator (CIDR) (Ruby)\e[0m"
    puts "Входные данные: #{@cidr}"
    puts
    puts "\e[32m─────────────────────────────────────────\e[0m"
    puts "IP-адрес:       #{info[:ip]}"
    puts "Маска (CIDR):   #{info[:cidr]} (#{info[:mask]})"
    puts "Сетевой адрес:  #{info[:network]}"
    puts "Широковещательный: #{info[:broadcast]}"
    puts "Диапазон хостов: #{info[:host_min]} — #{info[:host_max]}"
    puts "Количество хостов: #{info[:num_hosts]}"
    puts "Версия:         #{info[:version]}"
    puts "\e[32m─────────────────────────────────────────\e[0m"
  end

  def save_json(filename = 'ip_calc_output.json')
    File.write(filename, JSON.pretty_generate(info))
    puts "\e[32m💾 Сохранено: #{filename}\e[0m"
  end
end

def main
  if ARGV.length < 1
    puts "Usage: ruby ip_calculator.rb <IP/CIDR>"
    puts "Пример: ruby ip_calculator.rb 192.168.1.0/24"
    exit 1
  end

  cidr = ARGV[0]
  begin
    calc = IPCalculator.new(cidr)
    calc.print_info
    calc.save_json
  rescue => e
    puts "\e[31m❌ Ошибка: #{e.message}\e[0m"
    exit 1
  end
end

main if __FILE__ == $0
