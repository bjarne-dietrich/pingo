package utils

import (
	"fmt"
	"net"
	"os"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: ping\n")

	/*
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "    %v\n", f.Usage) // f.Name, f.Value
		})
	*/

}

func GetIPv4(addrs []string) net.IP {
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if IsIPv4(ip) {
			return ip
		}
	}

	return nil
}

func GetIPv6(addrs []string) net.IP {
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if IsIPv6(ip) {
			return ip
		}
	}

	return nil
}

func IsIPv4(ip net.IP) bool {
	return (ip.To4() != nil)
}

func IsIPv6(ip net.IP) bool {
	return (ip.To16() != nil)
}
