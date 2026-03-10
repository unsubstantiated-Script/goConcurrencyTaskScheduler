package scheduler

type PriorityQueue []*Task

// Len returns the number of elements in the PriorityQueue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less returns true if the priority of the first element is greater than the priority of the second element.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

// Swap swaps the elements at the given indices in the PriorityQueue.
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

// Push adds a new element to the PriorityQueue by appending it to the underlying slice.
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}

// Pop removes the last element from the PriorityQueue and returns it.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
