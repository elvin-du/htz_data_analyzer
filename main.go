package main

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
    Parse()

	addr := ":60003"
	log.Println("Runing at", addr)
	Run(addr)
}
