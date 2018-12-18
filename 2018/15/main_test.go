package main

import "testing"

func TestExample(t *testing.T) {
	observed, _ := solve("input_test", 3)
	if expected := 27140; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample2(t *testing.T) {
    observed, _ := solve("input_test2", 3)
    if expected := 27828; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample3(t *testing.T) {
	observed, _ := solve("input_combat1", 3)
    if expected := 36334; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestExample4(t *testing.T) {
	observed, _ := solve("input_combat2", 3)
	if expected := 39514; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}

func TestSolutionPart1(t *testing.T) {
	observed, _ := solve("input", 3)
	if expected := 269430; observed != expected {
		t.Errorf("Expected value was %d, got: %d.\n", expected, observed)
	}
}
