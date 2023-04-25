package input

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type InputReader interface {
	PromptAndRead(message string) (string, error)
	PromptAndReadPassword(message string) (string, error)
}

type DefaultInputReader struct {
	reader io.Reader
	writer io.Writer
}

func (ir *DefaultInputReader) PromptAndReadCredentials() (schoolId string, userId string, userPw string, err error) {
	fmt.Println("please confirm id/pw: ")
	schoolId, err = ir.PromptAndRead("school id: ")
	if err != nil {
		return "", "", "", fmt.Errorf("Error reading SCHOOL ID: %w", err)
	}
	userId, err = ir.PromptAndRead("user id: ")
	if err != nil {
		return "", "", "", fmt.Errorf("Error reading USER ID: %w", err)
	}
	userPw, err = ir.PromptAndReadPassword("user pw: ")
	if err != nil {
		return "", "", "", fmt.Errorf("Error reading USER PW: %w", err)
	}
	return schoolId, userId, userPw, err
}

func NewDefaultInputReader() *DefaultInputReader {
	return &DefaultInputReader{
		reader: os.Stdin,
		writer: os.Stdout,
	}
}

func PromptAndReadWithRW(reader io.Reader, writer io.Writer, message string) (string, error) {
	fmt.Fprint(writer, message)
	scanner := bufio.NewScanner(reader)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

func (ir *DefaultInputReader) PromptAndRead(message string) (string, error) {
	fmt.Fprint(ir.writer, message)
	scanner := bufio.NewScanner(ir.reader)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

func (ir *DefaultInputReader) PromptAndReadPassword(message string) (string, error) {
	fmt.Fprint(ir.writer, message)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Fprintln(ir.writer)
	return string(bytePassword), nil
}
