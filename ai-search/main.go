package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	DFS = iota
	BFS
	GBFS
	ASTAR
	DIJKSTRA
)

type Point struct {
	Row int
	Col int
}

type Wall struct {
	State Point
	Wall  bool
}

type Node struct {
	index  int
	State  Point
	Parent *Node
	Action string
}

type Solution struct {
	Actions []string
	Cells   []Point
}

type Maze struct {
	Height      int
	Width       int
	Start       Point
	Goal        Point
	Walls       [][]Wall
	CurrentNode *Node
	Solution    Solution
	Explored    []Point
	Steps       int
	NumExplored int
	Debug       bool
	SearchType  int
}

func main() {
	var m Maze
	var maze, searchType string
	flag.StringVar(&maze, "file", "data/maze.txt", "maze file (default is data/maze.txt)")
	flag.StringVar(&searchType, "search", "dfs", "search type")
	flag.Parse()

	err := m.Load(maze)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startTime := time.Now()
	switch searchType {
	case "dfs":
		m.SearchType = DFS
		solveDFS(&m)
	default:
		fmt.Printf("Invalid search type: %s\n", searchType)
		os.Exit(1)
	}

	if len(m.Solution.Actions) > 0 { // maze solved
		fmt.Println("Solved!")
		// m.PrintMaze()
		fmt.Println("Solution is ", len(m.Solution.Actions), "steps.")
		fmt.Println("Time to solve:", time.Since(startTime))
		m.OutputImage()
	} else {
		fmt.Println("No solution found!")
	}
	fmt.Println("Explored:", len(m.Explored), "nodes.")
}

func solveDFS(m *Maze) {
	var s DepthFirstSearch
	s.Game = m
	fmt.Println("Goal is:", m.Goal)
	s.Solve()
}

func (g *Maze) PrintMaze() {
	for r, row := range g.Walls {
		for c, col := range row {
			if col.Wall {
				fmt.Printf("█")
			} else if g.Start.Row == col.State.Row && g.Start.Col == col.State.Col {
				fmt.Printf("S")
			} else if g.Goal.Row == col.State.Row && g.Goal.Col == col.State.Col {
				fmt.Printf("G")
			} else if g.InSolution(Point{Row: r, Col: c}) {
				fmt.Printf("٭")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println("")
	}
}

func (g *Maze) InSolution(p Point) bool {
	return inExplored(p, g.Solution.Cells)
}

func (g *Maze) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening %s: %w\n", fileName, err)
	}
	defer f.Close()

	var fileContents []string
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading file %s: %w", fileName, err)
		}
		fileContents = append(fileContents, line)
	}

	foundStart, foundEnd := false, false
	for _, line := range fileContents {
		if strings.Contains(line, "A") {
			foundStart = true
		}
		if strings.Contains(line, "B") {
			foundEnd = true
		}
	}
	if !foundStart {
		return errors.New("starting location not found")
	}
	if !foundEnd {
		return errors.New("ending location not found")
	}
	g.Height = len(fileContents)
	g.Width = len(fileContents[0])
	var rows [][]Wall
	for i, row := range fileContents {
		var cols []Wall
		for j, col := range row {
			currLetter := fmt.Sprintf("%c", col)
			var wall Wall
			wall.State.Row = i
			wall.State.Col = j
			switch currLetter {
			case "A":
				g.Start = Point{Row: i, Col: j}
				wall.Wall = false
			case "B":
				g.Goal = Point{Row: i, Col: j}
				wall.Wall = false
			case " ":
				wall.Wall = false
			case "#":
				wall.Wall = true
			default:
				continue
			}
			cols = append(cols, wall)
		}
		rows = append(rows, cols)
	}
	g.Walls = rows
	return nil
}
