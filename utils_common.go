package main

import (
	"fmt"
	"log"
	"os/exec"
)

// ConsumeCPU consumes a given number of millicores for the specified duration.
func ConsumeCPU(millicores int, durationSec int) {
	log.Printf("ConsumeCPU millicores: %v, durationSec: %v", millicores, durationSec)
	// creating new consume cpu process
	arg1 := fmt.Sprintf("-millicores=%d", millicores)
	arg2 := fmt.Sprintf("-duration-sec=%d", durationSec)
	consumeCPU := exec.Command(consumeCPUBinary, arg1, arg2)
	err := consumeCPU.Run()
	if err != nil {
		log.Printf("Error while consuming CPU: %v", err)
	}
}

// GetCurrentStatus prints out a no-op.
func GetCurrentStatus() {
	log.Printf("GetCurrentStatus")
	// not implemented
}
