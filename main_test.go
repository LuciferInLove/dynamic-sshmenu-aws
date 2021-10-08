package main

import (
	"strings"
	"testing"
)

func TestParseResult(t *testing.T) {
	result := `{"Number":1,"IP":"172.16.0.11","Name":"test-instance","Zone":"us-east-1a"}`
	expected := instance{
		Number: 1,
		IP:     "172.16.0.11",
		Name:   "test-instance",
		Zone:   "us-east-1a",
	}

	actual, err := parseInstance(result)

	if err != nil {
		t.Fatalf("\nUnexpected error:\n %v", err.Error())
	}

	if actual != expected {
		t.Fatalf("\nexpected: %v\nactual:  %v", expected, actual)
	}

}

func TestParseResultError(t *testing.T) {
	result := `{"Number":"1","IP":"172.16.0.11","Name":"test-instance","Zone":"us-east-1a"}`
	errText := "json: cannot unmarshal string into Go struct field instance.Number of type int"
	_, err := parseInstance(result)

	if !strings.Contains(err.Error(), errText) {
		t.Fatalf("\nUnexpected error:\n %v", err.Error())
	}

}
