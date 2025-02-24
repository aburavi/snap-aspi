package main

import "testing"

func TestFoo(t *testing.T) {
    want := "Hello, World"
    if got := Foo(); got != want {
        t.Errorf("Foo() = %q, want %q", got, want)
    }
}
