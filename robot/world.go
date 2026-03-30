package main

import "time"

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
	return &Room{}
}
