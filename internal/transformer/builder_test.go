package transformer

import "testing"

func TestBasicReversibleStringBuffer(t *testing.T) {
	rsb := ReversibleStringBuilder{}
	rsb.Init()

	rsb.WriteString("hello ")
	rsb.WriteString("world")

	if rsb.Len() != 11 {
		t.Errorf("Expected %d\n Got %d\n", 11, rsb.Len())
	}

	if rsb.String() != "hello world" {
		t.Errorf("Expected %s\n Got %s\n", "hello world", rsb.String())
	}
}

func TestUndoReversibleStringBuffer(t *testing.T) {
	rsb := ReversibleStringBuilder{}
	rsb.Init()

	rsb.WriteString("hello ")
	rsb.WriteString("world. ")
	rsb.WriteString("how ")
	rsb.WriteString("are ")
	rsb.WriteString("you?")

	rsb.Reverse(3)

	if rsb.Len() != 13 {
		t.Errorf("Expected %d\n Got %d\n", 13, rsb.Len())
	}

	if rsb.String() != "hello world. " {
		t.Errorf("Expected %s\n Got %s\n", "hello world", rsb.String())
	}
}

func TestClearReversibleStringBuffer(t *testing.T) {
	rsb := ReversibleStringBuilder{}
	rsb.Init()

	rsb.WriteString("hello ")
	rsb.WriteString("world")

	rsb.Reverse(-1)

	if rsb.Len() != 0 {
		t.Errorf("Expected %d\n Got %d\n", 0, rsb.Len())
	}

	if rsb.String() != "" {
		t.Errorf("Expected %s\n Got %s\n", "hello world", rsb.String())
	}
}

func TestLenReversibleStringBuffer(t *testing.T) {
	rsb := ReversibleStringBuilder{}
	rsb.Init()

	rsb.WriteString("hello ")
	rsb.WriteString("world ")
	rsb.Flush()

	rsb.WriteString("more text")

	if rsb.Len() != 21 {
		t.Errorf("Expected %d\n Got %d\n", 21, rsb.Len())
	}
}
