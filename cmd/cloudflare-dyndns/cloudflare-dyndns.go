package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	cf "github.com/ole-treichel/cloudflare-dyndns/internal/cloudflare"
)

func getDynDns(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	zoneId := r.URL.Query().Get("zone_id")
	ip := r.URL.Query().Get("ip")
	domains := r.URL.Query()["domain"]

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: token param is required")
		return
	}

	if zoneId == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: zone_id param is required")
		return
	}

	if ip == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: ip param is required")
		return
	}

	if len(domains) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error: domain param is required")
		return
	}

	client := cf.NewClient(token)

	resp, err := client.GetDnsRecords(zoneId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Internal Server Error")
		return
	}

	for _, record := range resp.Result {
		if slices.Contains(domains, record.Name) {
			if record.Content != ip {
				fmt.Println(fmt.Sprintf("%s is out of date. old value: %s, new value: %s", record.Name, record.Content, ip))
				_, err := client.UpdateDnsRecord(zoneId, record.Id, record.Name, ip)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, "Internal Server Error")
					return
				} else {
					fmt.Println(fmt.Sprintf("%s updated to %s", record.Name, ip))
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Success")
	return
}

func main() {
	http.HandleFunc("/dyndns", getDynDns)

	fmt.Printf("Server listening on http://0.0.0.0:8000\n")
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed")
	} else {
		fmt.Printf("Error starting server %s\n", err)
		os.Exit(1)
	}
}
