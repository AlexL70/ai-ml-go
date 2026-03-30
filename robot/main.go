package main

import (
	"flag"
	"fmt"
)

func main() {
	var configFile, algorithm string
	var animate bool
	flag.StringVar(&configFile, "file", "data/empty.json", "Path to the configuration file")
	flag.StringVar(&algorithm, "algorithm", "random", "Cleaning algorithm")
	flag.BoolVar(&animate, "animate", true, "Enable animation while cleaning")
	flag.Parse()

	room := NewRoom(configFile, animate)
	fmt.Println(room)
}
