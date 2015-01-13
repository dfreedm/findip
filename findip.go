package main

import (
	"flag"
	"fmt"
	"net"
)

// return address string, is ip4, is external address
func processAddr(addr net.Addr) (string, bool, bool) {
	n, ok := addr.(*net.IPNet)
	if !ok {
		return "", false, false
	}
	ip := n.IP
	return ip.String(), ip.DefaultMask() != nil, ip.IsGlobalUnicast()
}

func abort(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		name    string
		version int
	)

	flag.IntVar(&version, "t", 0, "ipv4 or ipv6")
	flag.StringVar(&name, "n", "", "named interface")
	flag.Parse()

	ifaces, e := net.Interfaces()
	abort(e)

	for _, iface := range ifaces {

		if name != "" && iface.Name != name {
			continue
		}

		addrs, e := iface.Addrs()
		abort(e)

		for _, addr := range addrs {
			straddr, ip4, external := processAddr(addr)

			if !external {
				continue
			}
			if !ip4 && version == 4 {
				continue
			}
			if ip4 && version == 6 {
				continue
			}

			if name == "" {
				fmt.Printf("%s: ", iface.Name)
			}
			fmt.Println(straddr)
		}
	}
}
