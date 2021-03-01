package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Mkdir("test", 0777) // set up a temporary dir for generate files

	// Create whatever testfiles are needed in test/

	// Run all tests and clean up
	exitcode := m.Run()
	os.RemoveAll("test") // remove the directory and its contents.
	os.Exit(exitcode)
}

func TestTest(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
		failure  bool
	}{
		{"foobar", "foobar", false},
		{"a", "a", false},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			err := error(nil)
			actual := tt.input
			if tt.failure != true {
				if err != nil {
					t.Errorf("No error expected, got %v", err)
				}
				if tt.expected != "" {
					assert.Equal(t, tt.expected, actual, "should be equal")
				}
			} else {
				if err == nil {
					t.Errorf("Expected to fail?")
				}
			}
		})
	}
}
