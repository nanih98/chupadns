package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"text/tabwriter"
	"github.com/akamensky/argparse"
	"github.com/miekg/dns"
	"encoding/json"
	//"github.com/sirupsen/logrus"
)

func GetDNSServers(domain string) []string {
	var nameservers []string
	nameserver, _ := net.LookupNS(domain)
	for _, ns := range nameserver {
		nameservers = append(nameservers, ns.Host)
	}
	return nameservers
}

func GetDNSIps(nameservers [] string) []string {
	var ips []string
	for _, server := range nameservers {
		ip, err := net.LookupIP(server)
		if err != nil {
			log.Println("Error...")
		} 
		for _,value := range ip {
			ips = append(ips, value.String())
		}
	}
	return ips
}


func GetDNSAXFR(domain string, server string) {
	t := new(dns.Transfer)
	m := new(dns.Msg)

	m.SetAxfr(domain + ".")

	ch, err := t.In(m, server+":53")

	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)

	for env := range ch {
		data, err := json.MarshalIndent(env.RR, "", "  ")
		if err != nil {
			log.Fatalf("JSON Marshalling failed: %s", err)
		}
		fmt.Printf("%s\n",data)
		// if len(env.RR) > 0 {
		// 	fmt.Fprintln(w, "Go type\tName\tTTL\tClass\tRR type\tetc")
		// }
		// for _,result := range env.RR {
		// 	fmt.Fprintf(w, "%T\t%[1]s\n", result)
		// }
	}

	w.Flush()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	//Parse
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "Domain to scan"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	// contextLogger := logrus.WithFields(logrus.Fields{
	// 	"domain":  *domain,
	// 	"example": "test",
	// })

	nameservers := GetDNSServers(*domain)
	ips := GetDNSIps(nameservers)

	// Get ASXF opened zones
	for _,ip := range ips {
		GetDNSAXFR(*domain, ip)
	} 
}