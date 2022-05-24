package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"text/tabwriter"
	"github.com/akamensky/argparse"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

func GetDNSServers(domain string) []string {
	var nameservers []string
	nameserver, _ := net.LookupNS(domain)
	for _, ns := range nameserver {
		nameservers = append(nameservers, ns.Host)
	}
	return nameservers
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

	//fmt.Fprintln(w, "Go type\tName\tTTL\tClass\tRR type\tetc")

	for test := range ch {
		for _,result := range test.RR {
			fmt.Fprintf(w, "%T\t%[1]s\n", result)
		}
	}

	// for env := range ch {
	// 	if env.Error != nil {
	// 		err = env.Error
	// 		break
	// 	}
	// 	for _, rr := range env.RR {
	// 		fmt.Fprintf(w, "%T\t%[1]s\n", rr)
	// 	}
	// }

	w.Flush()

	fmt.Println("----------------------------------------------------------------")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//logrus.Out = os.Stdout
	//Parse
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "Domain to scan"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	contextLogger := logrus.WithFields(logrus.Fields{
		"domain":  *domain,
		"example": "test",
	})

	contextLogger.Info("I'll be logged with common and other field")

	nameservers := GetDNSServers(*domain)

	for _, server := range nameservers {
		ips, err := net.LookupIP(server)
		if err != nil {
			panic(err)
		}
		for _, ip := range ips {
			contextLogger.Info("Trying to solve %s in %s",*domain,ip)
			GetDNSAXFR(*domain, ip.String())
		}
	}
}
