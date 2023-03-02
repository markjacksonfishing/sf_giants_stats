package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// Set up a fake input string and a buffer to capture output
	inputStr := "SF\n"
	expectedOutput := "Year\tTm\t#Bat\tBatAge\tR/G\tG\tPA\tAB\tR\tH\t2B\t3B\tHR\tRBI\tSB\tCS\tBB\tSO\tBA\tOBP\tSLG\tOPS\tOPS+\n2022\tSFG\t27\t28.2\t4.4\t162\t6428\t5557\t708\t1431\t282\t43\t209\t694\t57\t20\t621\t1426\t.258\t.330\t.432\t.762\t106\n"
	outputBuf := bytes.NewBuffer(nil)

	// Run the function with the fake input and capture the output
	run(strings.NewReader(inputStr), outputBuf)
	actualOutput := outputBuf.String()

	// Compare the actual and expected output
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}
}
