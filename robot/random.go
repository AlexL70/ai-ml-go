package main

import (
	"fmt"
	"time"
)

func CleanRoomRandomWalk(room *Room, robot *Robot) {
	startTime := time.Now()
	moveCount := 0

	// Set variables
	maxMoves := room.Width * room.Height * 5 // Arbitrary large number to prevent infinite loops
	stuckCount := 0
	maxStuckCount := 5 // Maximum number of failed moves before changing strategy

	// Clean the current cell
	Clean(robot, room)
	if room.Animate {
		room.Display(robot, false)
		time.Sleep(moveDelay)
	}

	fmt.Println(stuckCount, maxStuckCount)
	for moveCount < maxMoves && room.CleanedCellCount < room.CleanableCellCount {
		// Generate a random angle (in radians)
		// Calculate a direction vector based on the angle
		// Use Bresenham's line algorithm to move in that direction until hitting an obstacle
		// If we did not move very much, increment the stuck counter and possibly change strategy
		// If stuck to many times, use A* to move to the nearest dirty cell
		// Add some adaptive behavior. Scan for dirty cells every ones in awhile.
	} // end for loop
	// Final sweep to ensure complete coverage

	cleaningTime := time.Since(startTime)
	displaySummary(room, robot, moveCount, cleaningTime)
}
