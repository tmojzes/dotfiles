package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseConsoleOutput(t *testing.T) {
	console := strings.Join([]string{
		"Started by user jenkins",
		"Found kubeconfig in prod-syd01 secret",
		"11:22:33 CUSTOM KUBX KUBECTL COMMAND: get nodes",
		"CUSTOM KUBX KUBECTL:",
		"11:22:33 NAME    STATUS",
		"11:22:34 node1   Ready",
		"",
		"--------------------",
		"11:22:35 CUSTOM KUBX KUBECTL COMMAND: get pods -A",
		"11:22:36 NS    POD",
		"",
		"====================",
		"Finished: SUCCESS",
	}, "\n")

	got := parseConsoleOutput(console)
	want := []commandBlock{
		{Command: "get nodes", Output: "NAME    STATUS\nnode1   Ready"},
		{Command: "get pods -A", Output: "NS    POD"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseConsoleOutput = %#v, want %#v", got, want)
	}
}

func TestParseConsoleOutputNoMarkers(t *testing.T) {
	console := "Started by user jenkins\nNothing to see here\n"
	if got := parseConsoleOutput(console); len(got) != 0 {
		t.Fatalf("parseConsoleOutput = %#v, want no blocks", got)
	}
}

func TestPairResultsMatching(t *testing.T) {
	requested := []string{"get nodes", "get pods -A"}
	parsed := []commandBlock{
		{Command: "get nodes", Output: "NAME STATUS"},
		{Command: "get pods -A", Output: "NS POD"},
	}
	results, warnings, err := pairResultsWithRequestedCommands(requested, parsed)
	if err != nil {
		t.Fatalf("pairResultsWithRequestedCommands: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}
	want := []commandResult{
		{Command: "get nodes", Output: "NAME STATUS"},
		{Command: "get pods -A", Output: "NS POD"},
	}
	if !reflect.DeepEqual(results, want) {
		t.Fatalf("results = %#v, want %#v", results, want)
	}
}

func TestPairResultsMismatch(t *testing.T) {
	requested := []string{"get nodes"}
	parsed := []commandBlock{{Command: "get pods", Output: "x"}}
	_, warnings, err := pairResultsWithRequestedCommands(requested, parsed)
	if err != nil {
		t.Fatalf("pairResultsWithRequestedCommands: %v", err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "command mismatch at position 1") {
		t.Fatalf("warnings = %#v, want one command mismatch warning", warnings)
	}
}

func TestPairResultsCountMismatch(t *testing.T) {
	requested := []string{"get nodes"}
	parsed := []commandBlock{
		{Command: "get nodes", Output: "x"},
		{Command: "get pods", Output: "y"},
	}
	results, warnings, err := pairResultsWithRequestedCommands(requested, parsed)
	if err != nil {
		t.Fatalf("pairResultsWithRequestedCommands: %v", err)
	}
	if len(warnings) != 2 {
		t.Fatalf("warnings = %#v, want count mismatch and extra block warnings", warnings)
	}
	if len(results) != 2 || results[1].Command != "get pods" {
		t.Fatalf("results = %#v, want the extra block appended", results)
	}
}

func TestPairResultsMissingBlock(t *testing.T) {
	requested := []string{"get nodes", "get pods"}
	parsed := []commandBlock{{Command: "get nodes", Output: "x"}}
	results, warnings, err := pairResultsWithRequestedCommands(requested, parsed)
	if err != nil {
		t.Fatalf("pairResultsWithRequestedCommands: %v", err)
	}
	if len(warnings) != 1 {
		t.Fatalf("warnings = %#v, want count mismatch warning", warnings)
	}
	if results[1].Output != "" || results[1].Command != "get pods" {
		t.Fatalf("results = %#v, want empty output for the missing block", results)
	}
}
