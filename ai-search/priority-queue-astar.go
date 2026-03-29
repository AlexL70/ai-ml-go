package main

type PriorityQueueAstar []*Node

func (pq PriorityQueueAstar) Len() int {
	return len(pq)
}

func (pq PriorityQueueAstar) Less(i, j int) bool {
	return pq[i].EstimatedCostToGoal < pq[j].EstimatedCostToGoal
}

func (pq PriorityQueueAstar) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueAstar) Push(n any) {
	l := len(*pq)
	item := n.(*Node)
	item.index = l
	*pq = append(*pq, item)
}

func (pq *PriorityQueueAstar) Pop() any {
	old := *pq
	l := len(old)
	item := old[l-1]
	old[l-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : l-1]
	return item
}
