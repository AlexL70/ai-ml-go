package main

import (
	"fmt"
	"math"
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

func bresenhamLine(x0, y0, x1, y1 int) []Point {
	// Initialize a slice to store all points on the line
	var points []Point

	// Calculate the absolute difference between the start and end points
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	// Determine the direction of movement along each axis
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	// Calculate the initial error value
	err := dx - dy

	for {
		// Append the current point to the result
		points = append(points, Point{X: x0, Y: y0})
		// Check to see if we've reached the end point
		if x0 == x1 && y0 == y1 {
			break
		}
		// Calculate the error for the next step
		err2 := err * 2
		// If moving in the x direction would keep us closer to the ideal line
		if err2 > -dy {
			// Check to see if we've reached the end point
			if x0 == x1 {
				break
			}
			// Update the error and move in the x direction
			err -= dy
			x0 += sx
		}
		// If moving in the y direction would keep us closer to the ideal line
		if err2 < dx {
			// Check to see if we've reached the end point
			if y0 == y1 {
				break
			}
			// Update the error and move in the y direction
			err += dx
			y0 += sy
		}
	}
	return points
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findNearestDirtyCell(room *Room, position Point) Point {
	var nearestCell Point = Point{X: -1, Y: -1}
	minDistance := math.MaxFloat64

	for i := range room.Width {
		for j := range room.Height {
			distance := heuristic(position, Point{X: i, Y: j})
			if distance < minDistance && !room.Grid[i][j].Cleaned && !room.Grid[i][j].Obstacle {
				minDistance = distance
				nearestCell = Point{X: i, Y: j}
			}
		}
	}
	return nearestCell
}
