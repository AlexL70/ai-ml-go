package main

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

const FloodedCost = 100

type AStar struct {
	Frontier PriorityQueueAstar
	Game     *Maze
}

func (as AStar) GetFrontier() []*Node {
	return as.Frontier
}

func (as *AStar) Add(i *Node) {
	i.CostToGoal = i.ManhattanDistance(as.Game.Start)
	i.EstimatedCostToGoal = euclideanDistance(i.State, as.Game.Goal) + float64(i.CostToGoal)
	if i.State.Water {
		i.EstimatedCostToGoal += FloodedCost
	}
	as.Frontier.Push(i)
	heap.Init(&as.Frontier)
}

func (as *AStar) ContainsState(i *Node) bool {
	for _, n := range as.Frontier {
		if n.State == i.State {
			return true
		}
	}
	return false
}

func (as *AStar) Empty() bool {
	return len(as.Frontier) == 0
}

func (as *AStar) Remove() (*Node, error) {
	if as.Empty() {
		return nil, errors.New("the frontier is empty")
	} else {
		if as.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range as.Frontier {
				fmt.Println("Node: ", x.State)
			}
		}
		return heap.Pop(&as.Frontier).(*Node), nil
	}
}

func (as *AStar) Solve() {
	fmt.Println("Starting to solve maze using A* search...")
	as.Game.NumExplored = 0

	start := Node{
		State:  as.Game.Start,
		Parent: nil,
		Action: "",
	}

	as.Add(&start)
	as.Game.CurrentNode = &start

	for {
		if as.Empty() {
			return // Cannot solve the maze
		}

		currentNode, err := as.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if as.Game.Debug {
			fmt.Println("Removed: ", currentNode.State)
			fmt.Println("----------")
			fmt.Println()
		}

		as.Game.CurrentNode = currentNode
		as.Game.NumExplored++
		as.Game.Explored = append(as.Game.Explored, currentNode.State)

		// Build animation frame if appropriate
		if as.Game.Animate {
			as.Game.OutputImage(fmt.Sprintf("./tmp/%06d.png", as.Game.NumExplored))
		}

		if currentNode.State == as.Game.Goal {
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
			as.Game.Solution = Solution{Actions: actions, Cells: cells}
			return
		}

		for _, x := range as.Neighbors(currentNode) {
			if !as.ContainsState(x) {
				if !inExplored(x.State, as.Game.Explored) {
					as.Add(&Node{State: x.State, Parent: currentNode, Action: x.Action})
				}
			}
		}
	}
}

func (as *AStar) Neighbors(node *Node) []*Node {
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
		if 0 <= n.State.Row && n.State.Row < as.Game.Height &&
			0 <= n.State.Col && n.State.Col < as.Game.Width &&
			!as.Game.Walls[n.State.Row][n.State.Col].Wall {
			if as.Game.Walls[n.State.Row][n.State.Col].State.Water {
				n.State.Water = true
			}
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
