package main

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type DijkstraSearch struct {
	Frontier PriorityQueueDijkstra
	Game     *Maze
}

func (ds DijkstraSearch) GetFrontier() []*Node {
	return ds.Frontier
}

func (ds *DijkstraSearch) Add(i *Node) {
	i.CostToGoal = i.ManhattanDistance(ds.Game.Start)
	ds.Frontier.Push(i)
	heap.Init(&ds.Frontier)
}

func (ds *DijkstraSearch) ContainsState(i *Node) bool {
	for _, n := range ds.Frontier {
		if n.State == i.State {
			return true
		}
	}
	return false
}

func (ds *DijkstraSearch) Empty() bool {
	return len(ds.Frontier) == 0
}

func (ds *DijkstraSearch) Remove() (*Node, error) {
	if ds.Empty() {
		return nil, errors.New("the frontier is empty")
	} else {
		if ds.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range ds.Frontier {
				fmt.Println("Node: ", x.State)
			}
		}
		return heap.Pop(&ds.Frontier).(*Node), nil
	}
}

func (ds *DijkstraSearch) Solve() {
	fmt.Println("Starting to solve maze using Dijkstra search...")
	ds.Game.NumExplored = 0

	start := Node{
		State:  ds.Game.Start,
		Parent: nil,
		Action: "",
	}

	ds.Add(&start)
	ds.Game.CurrentNode = &start

	for {
		if ds.Empty() {
			return // Cannot solve the maze
		}

		currentNode, err := ds.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if ds.Game.Debug {
			fmt.Println("Removed: ", currentNode.State)
			fmt.Println("----------")
			fmt.Println()
		}

		ds.Game.CurrentNode = currentNode
		ds.Game.NumExplored++
		ds.Game.Explored = append(ds.Game.Explored, currentNode.State)

		// Build animation frame if appropriate
		if ds.Game.Animate {
			ds.Game.OutputImage(fmt.Sprintf("./tmp/%06d.png", ds.Game.NumExplored))
		}

		if currentNode.State == ds.Game.Goal {
			var actions []string
			var cells []Point
			for {
				if currentNode.Parent != nil {
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}
			slices.Reverse(actions)
			slices.Reverse(cells)
			ds.Game.Solution = Solution{Actions: actions, Cells: cells}
			return
		}

		for _, x := range ds.Neighbors(currentNode) {
			if !ds.ContainsState(x) {
				if !inExplored(x.State, ds.Game.Explored) {
					ds.Add(&Node{State: x.State, Parent: currentNode, Action: x.Action})
				}
			}
		}
	}
}

func (ds *DijkstraSearch) Neighbors(node *Node) []*Node {
	row := node.State.Row
	col := node.State.Col
	candidates := []*Node{
		{State: Point{Row: row - 1, Col: col}, Parent: node, Action: "up"},
		{State: Point{Row: row + 1, Col: col}, Parent: node, Action: "down"},
		{State: Point{Row: row, Col: col - 1}, Parent: node, Action: "left"},
		{State: Point{Row: row, Col: col + 1}, Parent: node, Action: "right"},
	}

	var neighbors []*Node
	for _, n := range candidates {
		if 0 <= n.State.Row && n.State.Row < ds.Game.Height &&
			0 <= n.State.Col && n.State.Col < ds.Game.Width &&
			!ds.Game.Walls[n.State.Row][n.State.Col].Wall {
			neighbors = append(neighbors, n)
		}
	}
	// randomize neighbors before returning
	if len(neighbors) > 1 {
		for i := range neighbors {
			j := rand.Intn(len(neighbors) - 1)
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		}
	}

	return neighbors
}
