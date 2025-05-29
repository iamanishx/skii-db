package engine

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Engine struct {
	data map[string]int64
	file *os.File
	mu   sync.Mutex
}

var keyValueSeparator = " "

func NewEngine() (*Engine, error) {
	file, err := os.OpenFile("data.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file data:", err)
		return nil, err
	}
	return &Engine{
		data: make(map[string]int64),
		file: file,
		mu:   sync.Mutex{},
	}, nil
}

func (e *Engine) Set(key, value string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	offset, err := e.file.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return err
	}

	_, err = e.file.WriteString(key + keyValueSeparator + value + "\n")
	if err != nil {
		fmt.Println("Error appending text:", err)
		return err
	}

	e.data[key] = offset
	fmt.Println("Key set:", key, "at offset", offset)
	return nil
}

func (e *Engine) Get(key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.data[key]; !ok {
		return "", fmt.Errorf("key not found")
	}

	_, err := e.file.Seek(e.data[key]+int64(len(key))+1, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return "", err
	}

	buffer := make([]byte, 1)
	var content []byte

	for {
		n, err := e.file.Read(buffer)
		if err != nil {
			fmt.Println("Error reading file:", err)
			break
		}

		if n == 0 {
			break
		}

		if buffer[0] == '\n' {
			break
		}

		content = append(content, buffer[0])
	}
	return string(content), nil
}
