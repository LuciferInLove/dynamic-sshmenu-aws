package main

import (
	"strings"
	"testing"
)

func TestParseResult(t *testing.T) {
	result := "{1 172.16.0.11 test-instance us-east-1a}"
	expected := instance{
		Number: 1,
		IP:     "172.16.0.11",
		Name:   "test-instance",
		Zone:   "us-east-1a",
	}

	actual, err := parseResult(result)

	if err != nil {
		t.Fatalf("\nUnexpected error:\n %v", err.Error())
	}

	if actual != expected {
		t.Fatalf("\nexpected: %v\nactual:  %v", expected, actual)
	}

}

func TestParseResultError(t *testing.T) {
	result := "{n 172.16.0.11 test-instance us-east-1a}"
	errText := "invalid syntax"
	expected := instance{
		Number: 0,
		IP:     "",
		Name:   "",
		Zone:   "",
	}

	parsedResult, err := parseResult(result)

	if parsedResult != expected {
		t.Fatalf("\nUnexpected behaviour: result must not contains data")
	}

	if !strings.Contains(err.Error(), errText) {
		t.Fatalf("\nUnexpected error:\n %v", err.Error())
	}

}
