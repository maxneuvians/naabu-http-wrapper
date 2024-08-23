package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

type ScanResult struct {
	Host  string
	Ports []int
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")

	if host == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Host parameter is required"))
		return
	}

	ports := r.URL.Query().Get("ports")
	timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))

	if err != nil {
		timeout = 1000
	}

	options := runner.Options{
		Host:     goflags.StringSlice{host},
		ScanType: "s",
		Timeout:  timeout,
		OnResult: func(hr *result.HostResult) {

			portNumbers := make([]int, 0)

			for _, p := range hr.Ports {
				portNumbers = append(portNumbers, p.Port)
			}

			// Convert the result to a JSON object
			scanResult := ScanResult{
				Host:  hr.Host,
				Ports: portNumbers,
			}

			jsonResult, err := json.Marshal(scanResult)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error converting result to JSON"))
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
	log.Printf("About to listen on %s. Go to http://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
