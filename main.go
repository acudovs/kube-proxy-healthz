package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// parseKnownArgs parses known command line arguments ignoring unknown arguments
func parseKnownArgs() string {
	healthzBindAddress := flag.String(
		"healthz-bind-address", "127.0.0.1:10256", "Address to bind the health check server",
	)
	knownFlagArgs := map[string]int{
		"--healthz-bind-address": 1,
		"--help":                 0,
	}
	knownArgs := make([]string, 0)
	for i := 0; i < len(os.Args); i++ {
		if n, ok := knownFlagArgs[os.Args[i]]; ok {
			for j := 0; i+j < len(os.Args) && j <= n; j++ {
				knownArgs = append(knownArgs, os.Args[i+j])
			}
			i += n
		}
	}
	_ = flag.CommandLine.Parse(knownArgs)
	return *healthzBindAddress
}

// createServer creates a health check HTTP server
func createServer(healthzBindAddress string) *http.Server {
	server := &http.Server{
		Addr:         healthzBindAddress,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = fmt.Fprintf(w, `{"lastUpdated": %q,"currentTime": %q}`, currentTime, currentTime)
	})
	return server
}

// main starts a health check server
func main() {
	healthzBindAddress := parseKnownArgs()
	server := createServer(healthzBindAddress)
	log.Printf("Starting health check server on %s\n", healthzBindAddress)
	log.Fatal(server.ListenAndServe())
}
