package main

import (
	"flag"
	"time"

	"bitbucket.org/bertimus9/systemstat"
)

func main() {
	flag.Parse()
	// convert millicores to percentage
	millicoresPct := float64(*millicores) / float64(10)
	duration := time.Duration(*durationSec) * time.Second
	start := time.Now()
	first := systemstat.GetProcCPUSample()
	for time.Since(start) < duration {
		cpu := systemstat.GetProcCPUAverage(first, systemstat.GetProcCPUSample(), systemstat.GetUptime().Uptime)
		if cpu.TotalPct < millicoresPct {
			doSomething()
		} else {
			time.Sleep(sleep)
		}
	}
}
