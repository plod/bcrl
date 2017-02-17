package main

import (
	"log"
	"runtime"
	"strconv"
)

func init() {
	version, err := strconv.Atoi(runtime.Version()[4:])
	if err != nil {
		log.Fatalln("Couln't reliably find golang version [comment this code block out if you are using more than go1.8]", err)
	}
	if version < 8 {
		log.Fatalln("Minimum version go1.8 required for this project")
	}
}

func main() {
	log.Println("Great you have a good enough version")
}
