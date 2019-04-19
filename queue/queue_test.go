package queue

import "fmt"

func ExampleQueue_Pop() {
	pop := Queue{1}
	pop.Pop()
}

func ExampleQueue_Push() {
	pop := Queue{1}
	pop.Push(2)

	fmt.Println(pop.Pop())
	// Output:
	// 1
}

func ExampleQueue_Isempty() {

}
