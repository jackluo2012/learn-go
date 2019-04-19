package queue

//An FIFO queue
type Queue []int

//Pushes the element into the queue
func (q *Queue) Push(v int) {
	*q = append(*q, v)
}

//Pops element
func (q *Queue) Pop() int {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *Queue) Isempty() bool {
	return len((*q)) == 0
}
