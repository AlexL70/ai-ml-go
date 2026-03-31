package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/inancgumus/screen"
)

const (
	// Display characters
	charRobot          = "🔴"
	charWall           = "🟦"
	charFurniture      = "🪑"
	charClean          = "🧼"
	charDirty          = "🟫"
	charPath           = "🟩"
	charCat            = "🐈"
	catStopProbability = 0.1
	catStopDuration    = 5
	moveDelay          = 50 * time.Millisecond
	cellSize           = 10
)

type Point struct {
	X, Y int
}

type Cell struct {
	Type         string // wall, furniture, clean, dirty, bike
	Cleaned      bool
	Obstacle     bool
	ObstacleName string
}

type Furniture struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type Room struct {
	Grid               [][]Cell
	Width              int
	Height             int
	CleanableCellCount int
	CleanedCellCount   int
	Animate            bool
}

type RoomConfig struct {
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Furniture []Furniture `json:"furniture"`
}

func NewRoom(configFile string, animate bool) *Room {
	// Load from the json config file
	roomConfig, err := LoadRoomConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Convert room dimensions to grid cells
	gridWidth := roomConfig.Width / cellSize
	gridHeight := roomConfig.Height / cellSize
	grid := make([][]Cell, gridWidth)
	for i := range grid {
		grid[i] = make([]Cell, gridHeight)
		for j := range grid[i] {
			grid[i][j] = Cell{Type: "dirty", Cleaned: false, Obstacle: false}
		}
	}
	// Add walls
	for i := range gridWidth {
		grid[i][0] = Cell{Type: "wall", Obstacle: true, ObstacleName: "wall"}
		grid[i][gridHeight-1] = Cell{Type: "wall", Obstacle: true, ObstacleName: "wall"}
	}
	for j := range gridHeight {
		grid[0][j] = Cell{Type: "wall", Obstacle: true, ObstacleName: "wall"}
		grid[gridWidth-1][j] = Cell{Type: "wall", Obstacle: true, ObstacleName: "wall"}
	}

	// Add furniture
	for _, f := range roomConfig.Furniture {
		x := f.X / cellSize
		y := f.Y / cellSize
		width := f.Width / cellSize
		height := f.Height / cellSize
		for i := y; i < y+height; i++ {
			for j := x; j < x+width; j++ {
				grid[i][j] = Cell{Type: "furniture", Obstacle: true, ObstacleName: f.Name}
			}
		}
	}

	// Count cleanable cells
	cleanableCellCount := 0
	for i := range grid {
		for j := range grid[i] {
			if !grid[i][j].Obstacle {
				cleanableCellCount++
			}
		}
	}

	return &Room{
		Grid:               grid,
		Width:              gridWidth,
		Height:             gridHeight,
		CleanableCellCount: cleanableCellCount,
		CleanedCellCount:   0,
		Animate:            animate,
	}
}

func LoadRoomConfig(fileName string) (*RoomConfig, error) {
	// Read JSON
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	// parse JSON
	var config RoomConfig
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &config, nil
}

func (r *Room) Display(robot *Robot, showPath bool) {
	// Clear the console
	// fmt.Print("\033[H\033[2J") // Works on Unix-based systems only
	screen.Clear()
	for i := range r.Height {
		for j := range r.Width {
			cell := r.Grid[i][j]
			if robot.Position.X == i && robot.Position.Y == j {
				fmt.Print(charRobot)
			} else if showPath && isInPath(Point{X: i, Y: j}, robot.Path) {
				fmt.Print(charPath)
			} else {
				switch cell.Type {
				case "wall":
					fmt.Print(charWall)
				case "furniture":
					fmt.Print(charFurniture)
				case "clean":
					fmt.Print(charClean)
				case "dirty":
					fmt.Print(charDirty)
				}
			}
		}
		fmt.Println()
	}
	// Display cleaning progress
	progress := float64(r.CleanedCellCount) / float64(r.CleanableCellCount) * 100
	fmt.Printf("Cleaning progress: %.2f%% (%d/%d cells)\n", progress, r.CleanedCellCount, r.CleanableCellCount)
}

func isInPath(point Point, path []Point) bool {
	for _, p := range path {
		if p.X == point.X && p.Y == point.Y {
			return true
		}
	}
	return false
}

func displaySummary(room *Room, robot *Robot, moveCount int, cleaningTime time.Duration) {
	// Display the final room state with the robot's path
	fmt.Println("\nFinal room state with robot's path:")
	room.Display(robot, true)
	fmt.Println("\n======== Cleaning Summary ========")
	fmt.Printf("Room size: %dx%d (%d cm x %d cm)\n", room.Width, room.Height, room.Width*cellSize, room.Height*cellSize)

	// Calculate and display the percentage of the room cleaned
	percentCleaned := float64(room.CleanedCellCount) / float64(room.CleanableCellCount) * 100
	fmt.Printf("Cleaned %d cells out of %d which is (%.2f%%)\n", room.CleanedCellCount, room.CleanableCellCount, percentCleaned)
	// Display time and moves taken
	fmt.Printf("Total moves: %d\n", moveCount)
	fmt.Printf("Total cleaning time: %s\n", cleaningTime)
	// Calculate efficiency (cells cleaned per move)
	if moveCount > 0 {
		efficiency := float64(room.CleanedCellCount) / float64(moveCount)
		fmt.Printf("Cleaning efficiency: %.2f cells per move\n", efficiency)
	}
	fmt.Println("==================================")
}
