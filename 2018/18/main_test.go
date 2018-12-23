package main

import "testing"

func TestExample(t *testing.T) {
	if solve("input_test", 10, true) != 1147 {
		t.Error("Wrong")
	}
	if solve("input", 10, false) != 675100 {
		t.Error("Wrong")
	}
	solve("input", 1000000000, false)
}
