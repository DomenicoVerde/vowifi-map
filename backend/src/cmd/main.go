package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"vowifi_scanner/internal/persistence"
	"vowifi_scanner/internal/scan"
)

func main() {
	csvPath := flag.String("csv", "../data/mcc-mnc.csv", "Path to the mcc-mnc.csv dataset")
	mmdbPath := flag.String("mmdb", "../data/GeoLite2-City.mmdb", "Path to the GeoLite2-City.mmdb database")
	outputPath := flag.String("output", "../data/results.json", "Path to the output JSON report")
	workers := flag.Int("workers", 50, "Number of concurrent operators scanned at once")
	dnsTimeout := flag.Duration("dns-timeout", 3*time.Second, "Timeout for the EPDG DNS resolution")
	ikeTimeout := flag.Duration("ike-timeout", 3*time.Second, "Timeout waiting for an IKE_SA_INIT response")
	flag.Parse()

	// The internal/ike/message package logs verbosely to the standard
	// logger on every Encode/Decode call; silence it so it doesn't drown
	// out scan progress across thousands of operators.
	log.SetOutput(io.Discard)

	progress := log.New(os.Stderr, "", log.LstdFlags)

	operators, err := persistence.LoadOperators(*csvPath)
	if err != nil {
		progress.Fatalf("Failed to load operators: %v", err)
	}
	progress.Printf("Loaded %d operators from %s", len(operators), *csvPath)

	geo, err := persistence.OpenGeoIP(*mmdbPath)
	if err != nil {
		progress.Fatalf("Failed to open GeoIP database: %v", err)
	}
	defer geo.Close()

	options := scan.ScanOptions{
		DNSTimeout: *dnsTimeout,
		IKETimeout: *ikeTimeout,
	}

	results := make([]scan.Mno, len(operators))

	jobs := make(chan int)
	var wg sync.WaitGroup

	var done int
	var mu sync.Mutex
	reportProgress := func() {
		mu.Lock()
		done++
		d := done
		mu.Unlock()
		if d%100 == 0 || d == len(operators) {
			progress.Printf("Scanned %d/%d operators", d, len(operators))
		}
	}

	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = scan.ScanOperator(operators[i], geo, options)
				reportProgress()
			}
		}()
	}

	for i := range operators {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		progress.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(io.MultiWriter(outputFile, os.Stdout))
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		progress.Fatalf("Failed to write JSON output: %v", err)
	}

	progress.Printf("Done. Report written to %s", *outputPath)
}
