package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/port"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

type ScanResult struct {
	Host  string
	Ports []*port.Port
}

func scanHandler(w http.ResponseWriter, r *http.Request) {

	host := r.URL.Query().Get("host")
	ports := r.URL.Query().Get("ports")

	fmt.Println("Scanning host: ", host, " with ports: ", ports)

	options := runner.Options{
		Host:     goflags.StringSlice{host},
		ScanType: "s",
		Timeout:  1000,
		OnResult: func(hr *result.HostResult) {
			// Convert the result to a JSON object
			scanResult := ScanResult{
				Host:  hr.Host,
				Ports: hr.Ports,
			}

			jsonResult, err := json.Marshal(scanResult)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResult)

		},
		Ports: ports,
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer naabuRunner.Close()

	ctx := context.Background()
	naabuRunner.RunEnumeration(ctx)

}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/", scanHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
