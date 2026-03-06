package main

import (
	"testing"
)

func TestOpenBrowserFunc(t *testing.T) {
	var called bool
	origOpen := openBrowserFunc
	openBrowserFunc = func(url string) error {
		called = true
		if url != "https://myapp.localhost" {
			t.Fatalf("unexpected url: %s", url)
		}
		return nil
	}
	t.Cleanup(func() { openBrowserFunc = origOpen })

	if err := openBrowserFunc("https://myapp.localhost"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("openBrowserFunc not called")
	}
}
