package main

import (
	"fmt"
	"net"
	"os"

	flag "github.com/ogier/pflag"
)

const (
	tokenEnvName = "DO_API_TOKEN"
)

var (
	token = flag.StringP("token", "t", tokenEnvName, "DigitialOcean APIv2 token")
	v4opt = flag.BoolP("v4", "4", false, "Update IPv4 A record")
	v6opt = flag.BoolP("v6", "6", false, "Update IPv6 AAAA record")
)

func errorOut(msg string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", msg, err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *token == tokenEnvName {
		*token = os.Getenv(tokenEnvName)
	}

	if !*v4opt && !*v6opt {
		// If neither v4 or v6 is specified
		// default to v4
		*v4opt = true
	}

	var ipv4 net.IP
	var ipv6 net.IP

	addrs, err := getLocalAddresses()
	if err != nil {
		errorOut("getting local addresses", err)
	}
	for _, a := range addrs {
		v4 := (a.To4() != nil)

		if ipv4 == nil && v4 {
			ipv4 = a
		}
		if ipv6 == nil && !v4 {
			ipv6 = a
		}
	}

	if ipv4 == nil {
		ipv4, err = getNATAddress()
		if err != nil {
			errorOut("asking for our external address", err)
		}
	}

	args := flag.Args()
	if len(args) < 1 {
		hostname, err := os.Hostname()
		if err != nil {
			errorOut("no hostnames found", err)
		}
		args = append(args, hostname)
	}

	if *v6opt && ipv6 == nil {
		errorOut("AAAA records enabled and no public IPv6 address found.", nil)
	}

	if *v4opt && ipv4 == nil {
		errorOut("A records enabled and no public IPv4 address found.", nil)
	}

	for _, record := range args {
		if *v4opt {
			err := dyndnsUpdate(*token, record, "A", ipv4)
			if err != nil {
				errorOut("updaing DO DNS record", err)
			}
		}
		if *v6opt {
			err := dyndnsUpdate(*token, record, "AAAA", ipv6)
			if err != nil {
				errorOut("updaing DO DNS record", err)
			}
		}
	}
}
