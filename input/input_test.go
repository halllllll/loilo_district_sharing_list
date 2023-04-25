package input_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/github.com/halllllll/loilo_district_sharing_list/input"
)

type TestInputReader struct {
	reader *strings.Reader
	writer *bytes.Buffer
}

func (tir *TestInputReader) PromptAndRead(message string) (string, error) {
	return input.PromptAndReadWithRW(tir.reader, tir.writer, message)
}

func TestPromptAndRead(t *testing.T) {
	inputData := "test input with spaces\n"
	r := strings.NewReader(inputData)
	w := &bytes.Buffer{}

	result, err := input.PromptAndReadWithRW(r, w, "Dummy prompt: ")

	if err != nil {
		t.Fatalf("Error in PromptAndRead: %v", err)
	}

	expected := strings.TrimSpace(inputData)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}
}

func TestPromptAndReadWithRW(t *testing.T) {
	inputData := "test input with spaces\n"
	r := strings.NewReader(inputData)
	w := &bytes.Buffer{}

	result, err := input.PromptAndReadWithRW(r, w, "Dummy prompt: ")

	if err != nil {
		t.Fatalf("Error in PromptAndReadWithRW: %v", err)
	}

	expected := strings.TrimSpace(inputData)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	expectedPrompt := "Dummy prompt: "
	if w.String() != expectedPrompt {
		t.Errorf("Expected prompt: %q, got: %q", expectedPrompt, w.String())
	}
}
