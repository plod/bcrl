package main

import (
	"log"
	"runtime"
	"strconv"
)

func main() {
	version, err := strconv.ParseFloat(runtime.Version()[2:], 32)
	if err != nil {
		log.Fatalln("Couln't reliably find golang version [comment this code block out if you are using more than go1.8]", err)
	}
	log.Println(version)
}
