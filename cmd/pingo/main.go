package main

import (
	"flag"
	"fmt"
	"internal/icmp"
	"internal/utils"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var flagAlias = map[string]string{
	"host":    "h",
	"count":   "c",
	"verbose": "v",
}

func main() {
	var hostFlag string

	var ipv6Flag bool
	var verboseFlag bool

	var countFlag uint

	var destinationIP net.IP
	var destinationHost string

	log.SetOutput(os.Stdout)

	flag.StringVar(&hostFlag, "host", "", "Host")
	flag.StringVar(&hostFlag, "h", "", "Host")

	flag.BoolVar(&verboseFlag, "verbose", false, "Verbose output.")
	flag.BoolVar(&verboseFlag, "v", false, "Verbose output.")

	flag.UintVar(&countFlag, "count", 0, "")
	flag.UintVar(&countFlag, "c", 0, "")

	flag.BoolVar(&ipv6Flag, "6", false, "Use IPv6.")
	flag.BoolVar(&ipv6Flag, "ipv6", false, "Use IPv6.")

	flag.Usage = utils.Usage

	flag.Parse()

	// Check for no host given
	if flag.NArg() == 0 && hostFlag == "" {
		utils.Usage()
		os.Exit(64)
	}

	hostFlag = flag.Arg(0)

	// Set Seed for random
	rand.Seed(time.Now().UnixNano())

	// Check whether host is already an IP Address
	if ip := net.ParseIP(hostFlag); ip != nil {
		destinationHost = hostFlag
		destinationIP = ip

		// Automatically detect v6 without flag
		if utils.IsIPv6(ip) {
			ipv6Flag = true
		}

	} else {
		// Resolve Host
		addrs, err := net.LookupHost(hostFlag)
		if err != nil {
			fmt.Printf("ping: cannot resolve %s: Unknown host\n", hostFlag)
			os.Exit(68)
		}

		var ip net.IP = nil

		if ipv6Flag {
			ip = utils.GetIPv6(addrs)
		} else {
			ip = utils.GetIPv4(addrs)
		}

		if ip == nil {
			if ipv6Flag {
				fmt.Printf("ping: host is IPv4 only but mode was set to IPv6\n")
			} else {
				fmt.Printf("ping: host is IPv6 only but mode was set to IPv4\n")
			}

			os.Exit(1)
		}

		destinationHost = hostFlag
		destinationIP = ip

	}

	fmt.Printf("PING %s (%s): 56 data bytes\n", destinationHost, destinationIP)
	_ = ping(destinationIP, countFlag, ipv6Flag)

}

func ping(dst net.IP, count uint, ipv6 bool) (result []float32) {
	var identifier uint16 = uint16(rand.Uint32() & 0xffff)
	var sequenceNumber uint16 = 0

	conn, err := net.ListenPacket("ip4:icmp", "")
	if err != nil {
		panic(err)
	}
	result = make([]float32, 0)

	// Loop for multiple pings
	for i := 0; i < int(count); i++ {

		packet, err := icmp.NewICMPv4EchoRequestPacket(identifier, sequenceNumber, 56, []byte("Hello, Papa! "))
		if err != nil {
			panic(err)
		}
		reply := make([]byte, 256)

		start := time.Now()
		conn.WriteTo(packet.Raw(), &net.IPAddr{IP: dst})
		if err != nil {
			fmt.Errorf("writeTo failed: %v", err)
			return
		}

		// Wait for reply
		n, addr, err := conn.ReadFrom(reply)
		if err != nil {
			fmt.Errorf("ReadFrom failed: %v", err)
			return
		}
		// get ms
		elapsed := float32(time.Since(start).Microseconds()) / 1e3
		result = append(result, elapsed)

		fmt.Printf("%d bytes from %s: %.3f ms\n", n, addr, elapsed)

		time.Sleep(time.Second)

	}

	return result
}
