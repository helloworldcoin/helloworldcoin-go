package Stack

import "testing"

func TestStack(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(4)

	len := stack.Len()
	if len != 4 {
		t.Errorf("stack.Len() failed. Got %d, expected 4.", len)
	}

	value := stack.Peak().(int)
	if value != 4 {
		t.Errorf("stack.Peak() failed. Got %d, expected 4.", value)
	}

	value = stack.Pop().(int)
	if value != 4 {
		t.Errorf("stack.Pop() failed. Got %d, expected 4.", value)
	}

	len = stack.Len()
	if len != 3 {
		t.Errorf("stack.Len() failed. Got %d, expected 3.", len)
	}

	value = stack.Peak().(int)
	if value != 3 {
		t.Errorf("stack.Peak() failed. Got %d, expected 3.", value)
	}

	value = stack.Pop().(int)
	if value != 3 {
		t.Errorf("stack.Pop() failed. Got %d, expected 3.", value)
	}

	value = stack.Pop().(int)
	if value != 2 {
		t.Errorf("stack.Pop() failed. Got %d, expected 2.", value)
	}

	empty := stack.Empty()
	if empty {
		t.Errorf("stack.Empty() failed. Got %T, expected false.", empty)
	}

	value = stack.Pop().(int)
	if value != 1 {
		t.Errorf("stack.Pop() failed. Got %d, expected 1.", value)
	}

	empty = stack.Empty()
	if !empty {
		t.Errorf("stack.Empty() failed. Got %T, expected true.", empty)
	}

	nilValue := stack.Peak()
	if nilValue != nil {
		t.Errorf("stack.Peak() failed. Got %d, expected nil.", nilValue)
	}

	nilValue = stack.Pop()
	if nilValue != nil {
		t.Errorf("stack.Pop() failed. Got %d, expected nil.", nilValue)
	}

	len = stack.Len()
	if len != 0 {
		t.Errorf("stack.Len() failed. Got %d, expected 0.", len)
	}
}
