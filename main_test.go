package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// Set up a fake input string and a buffer to capture output
	inputStr := "SF\n"
	expectedOutput := "Predicted wins for the upcoming season: 97.31\n"
	outputBuf := bytes.NewBuffer(nil)

	// Run the function with the fake input and capture the output
	run(strings.NewReader(inputStr), outputBuf)
	actualOutput := outputBuf.String()

	// Compare the actual and expected output
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}

	// Test for invalid team abbreviation
	inputStr = "INVALID\n"
	expectedOutput = "Table not found in HTML\n"
	outputBuf.Reset()

	// Run the function with the fake input and capture the output
	run(strings.NewReader(inputStr), outputBuf)
	actualOutput = outputBuf.String()

	// Compare the actual and expected output
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}

	// Test for table with less than two columns
	inputStr = "SD\n"
	expectedOutput = "Error parsing data: row 0 has length 1 instead of 2\n"
	outputBuf.Reset()

	// Run the function with the fake input and capture the output
	run(strings.NewReader(inputStr), outputBuf)
	actualOutput = outputBuf.String()

	// Compare the actual and expected output
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}
}
