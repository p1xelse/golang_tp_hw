package main

type StackNode struct {
	data     interface{}
	nextNode *StackNode
}

type Stack struct {
	head *StackNode
}

func (stack *Stack) isEmpty() bool {
	return stack.head == nil
}

func (stack *Stack) Push(data interface{}) {
	stack.head = &StackNode{data, stack.head}
}

func (stack *Stack) Pop() (data interface{}) {
	if !stack.isEmpty() {
		data, stack.head = stack.head.data, stack.head.nextNode
	}

	return data
}

func (stack *Stack) Top() interface{} {
	if !stack.isEmpty() {
		return stack.head.data
	}

	return nil
}
