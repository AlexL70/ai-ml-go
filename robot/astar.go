package main

import (
	"container/heap"
	"math"
)

type PQItem struct {
	point    Point
	priority float64
	index    int
}

type PriorityQueue []*PQItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	l := len(*pq)
	item := x.(*PQItem)
	item.index = l
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	l := len(old)
	item := old[l-1]
	old[l-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : l-1]
	return item
}

func (pq *PriorityQueue) Update(item *PQItem, priority float64) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

func AStar(room *Room, start, goal Point) []Point {
	if !room.IsValid(start) || !room.IsValid(goal) {
		return []Point{}
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	openSetItems := make(map[Point]*PQItem)
	closedSet := make(map[Point]bool)
	gScore := make(map[Point]float64)
	fScore := make(map[Point]float64)
	cameFrom := make(map[Point]Point)

	gScore[start] = 0
	fScore[start] = heuristic(start, goal)
	startItem := &PQItem{point: start, priority: fScore[start], index: 0}
	heap.Push(pq, startItem)
	openSetItems[start] = startItem

	// Main A* loop
	for pq.Len() > 0 {
		// Get the node in the priority queue with the lowest fScore
		currentItem := heap.Pop(pq).(*PQItem)
		current := currentItem.point
		delete(openSetItems, current)

		// if we've reached the goal, reconstruct the path
		if current.X == goal.X && current.Y == goal.Y {
			return reconstructPath(cameFrom, current)
		}

		// Mark current point as processed
		closedSet[current] = true

		// Check neighbors
		for _, dir := range directions {
			neighbor := Point{X: current.X + dir[0], Y: current.Y + dir[1]}
			// Skip invalid and already processed neighbors
			if closedSet[neighbor] || !room.IsValid(neighbor) {
				continue
			}
			// Calculate tentative g-score
			tentativeGScore := gScore[current] + 1 // Assuming uniform cost for moving to a neighbor

			if _, exists := gScore[neighbor]; !exists || tentativeGScore < gScore[neighbor] {
				// Update the path information
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + heuristic(neighbor, goal)

				// Update the priority queue
				if item, exists := openSetItems[neighbor]; exists {
					pq.Update(item, fScore[neighbor])
				} else {
					neighborItem := &PQItem{point: neighbor, priority: fScore[neighbor]}
					heap.Push(pq, neighborItem)
					openSetItems[neighbor] = neighborItem
				}
			}
		}
	}
	// Goal was not reached
	return []Point{}
}

func heuristic(a, b Point) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

func reconstructPath(cameFrom map[Point]Point, current Point) []Point {
	var totalPath []Point
	totalPath = append(totalPath, current)
	for {
		prev, exists := cameFrom[current]
		if !exists {
			break
		}
		current = prev
		totalPath = append([]Point{current}, totalPath...)
	}
	return totalPath
}
