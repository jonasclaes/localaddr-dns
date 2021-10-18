package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/miekg/dns"
)

var base_domain string
var port, ttl int

func main() {
	flag.IntVar(&port, "port", 8090, "Port to run on")
	flag.IntVar(&ttl, "ttl", 3600, "Time to live")
	flag.StringVar(&base_domain, "base_domain", "localaddr.net", "Base domain")
	flag.Parse()

	dns.HandleFunc(".", handleRequest)

	// UDP server
	go func() {
		srv := &dns.Server{
			Addr: ":" + strconv.Itoa(port),
			Net:  "udp",
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set UDP listener %s\n", err.Error())
		}
	}()

	// TCP server
	go func() {
		srv := &dns.Server{
			Addr: ":" + strconv.Itoa(port),
			Net:  "tcp",
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set TCP listener %s\n", err.Error())
		}
	}()

	log.Printf("localaddr-dns server started up on port %d.\n", port)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name

	regexIPv4, _ := regexp.Compile("^([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.")
	regexClassA, _ := regexp.Compile("^(10)\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.")
	regexClassB, _ := regexp.Compile("^(172)\\.(1[6-9]|2[0-9]|3[0-2])\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.")
	regexClassC, _ := regexp.Compile("^(192)\\.(168)\\.([0-9][0-9]?[0-9]?)\\.([0-9][0-9]?[0-9]?)\\.")

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{
		Name:   domain,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}

	if strings.Contains(domain, base_domain) && (regexClassA.MatchString(domain) || regexClassB.MatchString(domain) || regexClassC.MatchString(domain)) {
		data := regexIPv4.FindStringSubmatch(domain)
		if ip := net.ParseIP(fmt.Sprintf("%s.%s.%s.%s", data[1], data[2], data[3], data[4])); ip != nil {
			rr.A = ip
		} else {
			rr.A = net.ParseIP("127.0.0.1")
		}
	} else {
		rr.A = net.ParseIP("127.0.0.1")
	}

	m.Answer = []dns.RR{rr}
	w.WriteMsg(m)
}