package main

import (
	"runtime"
	"strconv"
)

func main() {
	version, err := strconv.ParseFloat(runtime.Version()[2:])
}
