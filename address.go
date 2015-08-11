package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var (
	privateRanges = []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}
)

func getLocalAddresses() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	privates := make([]net.IPNet, 0, len(privateRanges))
	for _, s := range privateRanges {
		_, ipnet, err := net.ParseCIDR(s)
		if err != nil {
			panic(err)
		}
		privates = append(privates, *ipnet)
	}

	result := make([]net.IP, 0, 5)

addrLoop:
	for _, a := range addrs {
		var ip net.IP

		if ipaddr, ok := a.(*net.IPAddr); ok {
			ip = ipaddr.IP
		}

		if ipnet, ok := a.(*net.IPNet); ok {
			ip = ipnet.IP
		}

		if !ip.IsGlobalUnicast() {
			continue
		}

		for _, n := range privates {
			if n.Contains(ip) {
				continue addrLoop
			}
		}

		result = append(result, ip)
	}

	return result, nil
}

func getNATAddress() (net.IP, error) {
	resp, err := http.Get("http://ip4.telize.com")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(strings.TrimSpace(string(data)))
	return ip, nil
}
