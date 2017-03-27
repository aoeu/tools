package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {
	args := struct {
		hostname     string
		nameServerIP string
	}{}
	flag.StringVar(&args.hostname, "hostname", "example.com", "The hostname to fetch the IP address of via a DNS server.")
	flag.StringVar(&args.nameServerIP, "with", "8.8.8.8", "The IP address of a DNS server to fetch hostnames from.")
	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Unknown arguments were provided, exiting : %v\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(args.hostname+".", dns.TypeA)
	resp, _, err := c.Exchange(&m, args.nameServerIP+":53")
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error encountered when looking up host %v on DNS server %v : %v\n", args.hostname, args.nameServerIP, err)
		os.Exit(1)
	case len(resp.Answer) == 0:
		fmt.Fprintf(os.Stderr, "No IP addresses received in response when looking up host %v on DNS server %v\n", args.hostname, args.nameServerIP)
		os.Exit(1)
	}
	for _, ans := range resp.Answer {
		Arecord := ans.(*dns.A)
		fmt.Fprintf(os.Stdout, "%v\n", Arecord.A)
	}
}
