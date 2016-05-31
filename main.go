package main

import (
	"htz_data_analyzer/log"
)

//func init() {
//	log.SetFlags(log.Lshortfile)
//}

func main() {
	Parse()
	addr := ":60003"
	log.Infoln("Runing at", addr)
	Run(addr)
}
