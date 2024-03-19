package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run()

	// Testing if error is nil
	if err != nil {
		t.Error("Failed run")
	}
}
