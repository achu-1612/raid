package main

import (
	"log"

	"github.com/achu-1612/raid"
)

func main() {
	r, err := raid.New(raid.RAIDType10, "raid")
	if err != nil {
		log.Fatalf("failed to create RAID: %v", err)
	}

	r.Write("example.txt", "Hello RAID 0!")
	data, err := r.Read("example.txt")
	if err != nil {
		log.Fatalf("failed to read from RAID: %v", err)
	}

	log.Println("RAID 0 Read:", data)
}
