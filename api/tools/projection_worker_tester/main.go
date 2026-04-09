package main

import (
	"api/services"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	instance := flag.String("instance", "local", "Instance name for logs")
	interval := flag.Duration("interval", 5*time.Second, "Interval between claim attempts")
	runFor := flag.Duration("duration", 1*time.Minute, "Total duration to run the tester")
	flag.Parse()

	fmt.Printf("[tester %s] starting, interval=%s, duration=%s\n", *instance, *interval, *runFor)

	// Initialize DB and tables (no migrations unless needed)
	services.SetupTables(false, false)
	if services.DB == nil {
		fmt.Fprintf(os.Stderr, "[tester %s] DB not initialized. Check configuration.\n", *instance)
		os.Exit(1)
	}

	end := time.Now().Add(*runFor)
	for time.Now().Before(end) {
		services.ProcessQueuedProjectionJobs()
		time.Sleep(*interval)
	}

	fmt.Printf("[tester %s] done.\n", *instance)
}
