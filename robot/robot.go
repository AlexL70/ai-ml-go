package main

var directions = [][]int{
	{0, -1}, // Up
	{1, 0},  // Right
	{0, 1},  // Down
	{-1, 0}, // Left
}

type Robot struct {
	Position             Point
	Path                 []Point
	CleanRoom            func(*Room, *Robot)
	Direction            float64
	ObstaclesEncountered map[string]bool
}

func NewRobot(startX, startY int) *Robot {
	return &Robot{
		Position:             Point{X: startX, Y: startY},
		Path:                 []Point{{X: startX, Y: startY}},
		ObstaclesEncountered: make(map[string]bool),
	}
}

func Clean(robot *Robot, room *Room) {
	// Clean the current cell
	currentCell := &room.Grid[robot.Position.X][robot.Position.Y]
	if !currentCell.Cleaned && !currentCell.Obstacle {
		currentCell.Cleaned = true
		currentCell.Type = "clean"
		room.CleanedCellCount++
	}

	CheckAdjacentObstacles(robot, room)
}

func CheckAdjacentObstacles(robot *Robot, room *Room) {
	x, y := robot.Position.X, robot.Position.Y
	for _, dir := range directions {
		newX, newY := x+dir[0], y+dir[1]
		RecordObstacle(robot, room, newX, newY)
	}
}

func RecordObstacle(robot *Robot, room *Room, x, y int) {
	if x >= 0 && x < room.Width && y >= 0 && y < room.Height {
		cell := room.Grid[x][y]
		if cell.Obstacle && cell.Type == "furniture" && cell.ObstacleName != "" {
			robot.ObstaclesEncountered[cell.ObstacleName] = true
		}
	}
}
