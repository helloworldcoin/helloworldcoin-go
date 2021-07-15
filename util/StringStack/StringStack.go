package StringStack

import "container/list"

type StringStack struct {
	list *list.List
}

func NewStringStack() *StringStack {
	list := list.New()
	return &StringStack{list}
}

func (stack *StringStack) Push(value string) {
	stack.list.PushBack(value)
}

func (stack *StringStack) Pop() *string {
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		popValue := e.Value.(string)
		return &popValue
	}
	return nil
}

func (stack *StringStack) Peek() *string {
	e := stack.list.Back()
	if e != nil {
		peekValue := e.Value.(string)
		return &peekValue
	}

	return nil
}

func (stack *StringStack) Size() int {
	return stack.list.Len()
}

func (stack *StringStack) IsEmpty() bool {
	return stack.list.Len() == 0
}
