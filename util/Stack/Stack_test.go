package Stack

/*
 @author x.king xdotking@gmail.com
*/

import "testing"

func TestStack(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(4)

	len := stack.Size()
	if len != 4 {
		t.Errorf("stack.Len() failed. Got %d, expected 4.", len)
	}

	value := stack.Peek().(int)
	if value != 4 {
		t.Errorf("stack.Peak() failed. Got %d, expected 4.", value)
	}

	value = stack.Pop().(int)
	if value != 4 {
		t.Errorf("stack.Pop() failed. Got %d, expected 4.", value)
	}

	len = stack.Size()
	if len != 3 {
		t.Errorf("stack.Len() failed. Got %d, expected 3.", len)
	}

	value = stack.Peek().(int)
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

	empty := stack.IsEmpty()
	if empty {
		t.Errorf("stack.Empty() failed. Got %T, expected false.", empty)
	}

	value = stack.Pop().(int)
	if value != 1 {
		t.Errorf("stack.Pop() failed. Got %d, expected 1.", value)
	}

	empty = stack.IsEmpty()
	if !empty {
		t.Errorf("stack.Empty() failed. Got %T, expected true.", empty)
	}

	nilValue := stack.Peek()
	if nilValue != nil {
		t.Errorf("stack.Peak() failed. Got %d, expected nil.", nilValue)
	}

	nilValue = stack.Pop()
	if nilValue != nil {
		t.Errorf("stack.Pop() failed. Got %d, expected nil.", nilValue)
	}

	len = stack.Size()
	if len != 0 {
		t.Errorf("stack.Len() failed. Got %d, expected 0.", len)
	}
}
