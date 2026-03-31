package main

import (
	"math"
	"math/rand"
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

	for moveCount < maxMoves && room.CleanedCellCount < room.CleanableCellCount {
		// Generate a random angle (in radians)
		angle := rand.Float64() * 2 * math.Pi
		// Calculate a direction vector based on the angle
		dx := math.Cos(angle)
		dy := math.Sin(angle)
		// Use Bresenham's line algorithm to move in that direction until hitting an obstacle
		moves := moveAtAngleUntilObstacle(room, robot, dx, dy)
		moveCount += moves
		// If we did not move very much, increment the stuck counter and possibly change strategy
		if moves < 3 {
			stuckCount++
			// If stuck to many times, use A* to move to the nearest dirty cell
			if stuckCount >= maxStuckCount {
				stuckCount = 0
				nearestDirty := findNearestDirtyCell(room, robot.Position)
				if nearestDirty.X != -1 && nearestDirty.Y != -1 {
					path := AStar(room, robot.Position, nearestDirty)
					if len(path) > 0 {
						// Move along the A* path
						for _, point := range path {
							robot.Position = point
							robot.Path = append(robot.Path, robot.Position)
							Clean(robot, room)
							moveCount++
							if room.Animate {
								room.Display(robot, false)
								time.Sleep(moveDelay)
							}
						}
					}
				} else {
					stuckCount = 0 // Reset the stuck count to give the random walk more chances before switching to A*
				}
			}
		}
		// Add some adaptive behavior. Scan for dirty cells every ones in awhile.
		if moveCount%20 == 0 {
			if rand.Float64() < 0.3 {
				dirtyCell := findNearestDirtyCell(room, robot.Position)
				if dirtyCell.X != -1 && dirtyCell.Y != -1 {
					path := AStar(room, robot.Position, dirtyCell)
					// Move along the A* path
					if len(path) > 0 {
						for _, point := range path {
							robot.Position = point
							robot.Path = append(robot.Path, robot.Position)
							Clean(robot, room)
							moveCount++
							if room.Animate {
								room.Display(robot, false)
								time.Sleep(moveDelay)
							}
						}
					}
				}
			}
		}
	} // end for loop
	// Final sweep to ensure complete coverage
	for x := range room.Width {
		for y := range room.Height {
			if !room.Grid[x][y].Cleaned && !room.Grid[x][y].Obstacle {
				path := AStar(room, robot.Position, Point{X: x, Y: y})
				if len(path) == 0 {
					continue
				}
				for _, point := range path {
					robot.Position = point
					robot.Path = append(robot.Path, robot.Position)
					Clean(robot, room)
					moveCount++
					if room.Animate {
						room.Display(robot, false)
						time.Sleep(moveDelay)
					}
				}
			}
		}
	}
	cleaningTime := time.Since(startTime)
	displaySummary(room, robot, moveCount, cleaningTime)
}

func moveAtAngleUntilObstacle(room *Room, robot *Robot, dx, dy float64) int {
	moveCount := 0
	maxDistance := math.Max(float64(room.Width), float64(room.Height)) * 2 // Maximum distance to move
	startX, startY := robot.Position.X, robot.Position.Y
	endX := int(float64(startX) + dx*maxDistance)
	endY := int(float64(startY) + dy*maxDistance)

	points := bresenhamLine(startX, startY, endX, endY)
	// Move along the line until we hit an obstacle
	for i := 1; i < len(points); i++ {
		x, y := points[i].X, points[i].Y
		if !room.IsValid(Point{X: x, Y: y}) {
			break
		}
		robot.Position = Point{X: x, Y: y}
		robot.Path = append(robot.Path, robot.Position)
		Clean(robot, room)
		moveCount++

		if room.Animate {
			room.Display(robot, false)
			time.Sleep(moveDelay)
		}
	}
	return moveCount
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
