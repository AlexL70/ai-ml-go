package main

type PriorityQueueGBFS []*Node

func (pq PriorityQueueGBFS) Len() int {
	return len(pq)
}

func (pq PriorityQueueGBFS) Less(i, j int) bool {
	return pq[i].CostToGoal < pq[j].CostToGoal
}

func (pq PriorityQueueGBFS) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueGBFS) Push(n any) {
	l := len(*pq)
	item := n.(*Node)
	item.index = l
	*pq = append(*pq, item)
}

func (pq *PriorityQueueGBFS) Pop() any {
	old := *pq
	l := len(old)
	item := old[l-1]
	old[l-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : l-1]
	return item
}
