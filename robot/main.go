package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var configFile, algorithm string
	var animate bool
	flag.StringVar(&configFile, "file", "data/empty.json", "Path to the configuration file")
	flag.StringVar(&algorithm, "algorithm", "random", "Cleaning algorithm")
	flag.BoolVar(&animate, "animate", true, "Enable animation while cleaning")
	flag.Parse()

	room := NewRoom(configFile, animate)

	// Create robot
	robot := NewRobot(1, 1)

	// Assign cleaning algorithm
	switch algorithm {
	case "random":
		robot.CleanRoom = CleanRoomRandomWalk
	default:
		fmt.Printf("Unknown algorithm: %s\n", algorithm)
		os.Exit(1)
	}

	// Clean the room
	robot.CleanRoom(room, robot)
}
