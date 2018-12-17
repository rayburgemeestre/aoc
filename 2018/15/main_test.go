package main

import "testing"

func TestExample(t *testing.T) {
	if observed, expected := solve("input_test", 3), 27140; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample2(t *testing.T) {
	if observed, expected := solve("input_test2", 3), 27828; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample3(t *testing.T) {
	if observed, expected := solve("input_combat1", 3), 36334; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample4(t *testing.T) {
	if observed, expected := solve("input_combat2", 3), 39514; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestSolutionPart1(t *testing.T) {
	if observed, expected := solve("input", 3), 269430; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}
