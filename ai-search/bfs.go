package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type BreadthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

func (bfs BreadthFirstSearch) GetFrontier() []*Node {
	return bfs.Frontier
}

func (bfs *BreadthFirstSearch) Add(i *Node) {
	bfs.Frontier = append(bfs.Frontier, i)
}

func (bfs *BreadthFirstSearch) ContainsState(i *Node) bool {
	for _, n := range bfs.Frontier {
		if n.State == i.State {
			return true
		}
	}
	return false
}

func (bfs *BreadthFirstSearch) Empty() bool {
	return len(bfs.Frontier) == 0
}

func (bfs *BreadthFirstSearch) Remove() (*Node, error) {
	if bfs.Empty() {
		return nil, errors.New("the frontier is empty")
	} else {
		if bfs.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range bfs.Frontier {
				fmt.Println("Node: ", x.State)
			}
		}
		node := bfs.Frontier[0]
		bfs.Frontier = bfs.Frontier[1:]
		return node, nil
	}
}

func (bfs *BreadthFirstSearch) Solve() {
	fmt.Println("Starting to solve maze using Depth First Search...")
	bfs.Game.NumExplored = 0

	start := Node{
		State:  bfs.Game.Start,
		Parent: nil,
		Action: "",
	}

	bfs.Add(&start)
	bfs.Game.CurrentNode = &start

	for {
		if bfs.Empty() {
			return // Cannot solve the maze
		}

		currentNode, err := bfs.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if bfs.Game.Debug {
			fmt.Println("Removed: ", currentNode.State)
			fmt.Println("----------")
			fmt.Println()
		}

		bfs.Game.CurrentNode = currentNode
		bfs.Game.NumExplored++
		bfs.Game.Explored = append(bfs.Game.Explored, currentNode.State)

		// Build animation frame if appropriate
		if bfs.Game.Animate {
			bfs.Game.OutputImage(fmt.Sprintf("./tmp/%06d.png", bfs.Game.NumExplored))
		}

		if currentNode.State == bfs.Game.Goal {
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
			bfs.Game.Solution = Solution{Actions: actions, Cells: cells}
			return
		}

		for _, x := range bfs.Neighbors(currentNode) {
			if !bfs.ContainsState(x) {
				if !inExplored(x.State, bfs.Game.Explored) {
					bfs.Add(&Node{State: x.State, Parent: currentNode, Action: x.Action})
				}
			}
		}
	}
}

func (bfs *BreadthFirstSearch) Neighbors(node *Node) []*Node {
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
		if 0 <= n.State.Row && n.State.Row < bfs.Game.Height &&
			0 <= n.State.Col && n.State.Col < bfs.Game.Width &&
			!bfs.Game.Walls[n.State.Row][n.State.Col].Wall {
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
