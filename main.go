package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

type RemoteIP struct {
	RemoteIP string `json:"remote_ip"`
}

// REFERENCE: https://www.name.com/api-docs/DNS#ListRecords

/*
	curl -u 'username:token' 'https://api.dev.name.com/v4/domains/example.org/records/12345'
*/

type Record struct {
	ID         int32  `json:"id"`
	DomainName string `json:"string"`
	Host       string `json:"host"`
	FQDN       string `json:"fqdn"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	TTL        uint32 `json:"ttl"`
	Priority   uint32 `json:"priority"`
}

var (
	baseURL = "https://api.dev.name.com/"
)

type ListRecords struct {
	Records  []Record `json:"records"`
	NextPage int32    `json:"nextPage"`
	LastPage int32    `json:"lastPage"`
}

// GET /v4/domains/{domainName}/records
func listRecords(domainName string) ListRecords {
	resp, err := http.Get(fmt.Sprintf("%s/v4/domains/%s/records", baseURL, domainName))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		log.Fatalf("%d %s ", resp.StatusCode, resp.Status)
	}

	var model ListRecords
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		log.Fatal(err)
	}

	return model
}
func getRecord(id int32, domainName string) {

}

func compare(dnsHost, whatismyipURL string) bool {
	host, err := net.LookupHost(dnsHost)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(whatismyipURL)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	var remoteIp RemoteIP
	if err := json.NewDecoder(response.Body).Decode(&remoteIp); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hosts", host)
	fmt.Println("Response", remoteIp)

	for _, host := range host {
		if remoteIp.RemoteIP != host {
			return true
		}
	}
	return false
}

var (
	dnsHost       = flag.String("dns-host", "", "")
	whatismyipURL = flag.String("remote-url", "", "")
	domain        = flag.String("domain", "", "")
)

func main() {
	flag.Parse()
	compare(*dnsHost, *whatismyipURL)
	listRecords(*domain)
}
