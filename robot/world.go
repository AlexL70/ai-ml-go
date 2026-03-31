package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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

	// TODO: Add furniture

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
