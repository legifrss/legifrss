package main

import "testing"

func TestStart(t *testing.T) {
	want := "ok"
	if got := Start(); got != want {
		t.Errorf("Start() = %q, want %q", got, want)
	}
}
