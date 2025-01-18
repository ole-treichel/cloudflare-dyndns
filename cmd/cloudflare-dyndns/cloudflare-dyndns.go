package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"

	cf "github.com/ole-treichel/cloudflare-dyndns/internal/cloudflare"
	c "github.com/ole-treichel/cloudflare-dyndns/internal/config"
)

var config c.Config

func getDynDns(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	ip := r.URL.Query().Get("ip")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: token param is required")
		return
	}

	if ip == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: ip param is required")
		return
	}

	client := cf.NewClient(token)

	for _, domainConfig := range config.DomainConfigs {
		resp, err := client.GetDnsRecords(domainConfig.ZoneId)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Internal Server Error")
			return
		}

		for _, record := range resp.Result {
			if slices.Contains(domainConfig.Domains, record.Name) {
				if record.Content != ip {
					log.Println(fmt.Sprintf("%s is out of date. old value: %s, new value: %s", record.Name, record.Content, ip))
					_, err := client.UpdateDnsRecord(domainConfig.ZoneId, record.Id, record.Name, ip)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						io.WriteString(w, "Internal Server Error")
						return
					} else {
						log.Println(fmt.Sprintf("%s updated to %s", record.Name, ip))
					}
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Success")
	return
}

func main() {
	var err error
	config, err = c.GetConfig()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/dyndns", getDynDns)
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8000"
	}

	fmt.Printf("Server listening on http://0.0.0.0:%s\n", httpPort)
	err = http.ListenAndServe(":"+httpPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed")
	} else {
		fmt.Printf("Error starting server %s\n", err)
		os.Exit(1)
	}
}
