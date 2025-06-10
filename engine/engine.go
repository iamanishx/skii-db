package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	data       map[string]int64
	file       *os.File
	fileDelete *os.File
	mu         sync.Mutex
	muDlete    sync.Mutex
	fileMu     sync.Mutex
}

type Item struct {
	Key    string
	Value  string
	Offset int64
}

func (e *Engine) Lock() {
	panic("unimplemented")
}

var keyValueSeparator = " "

func NewEngine() (*Engine, error) {
	file, err := os.OpenFile("data.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file data:", err)
		return nil, err
	}
	fileDelete, err := os.OpenFile("delete.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file delete:", err)
		return nil, err
	}

	return &Engine{
		data:       make(map[string]int64),
		file:       file,
		fileDelete: fileDelete,
		mu:         sync.Mutex{},
		muDlete:    sync.Mutex{},
		fileMu:     sync.Mutex{},
	}, nil
}

func (e *Engine) Set(key, value string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if strings.Contains(key, " ") {
		return fmt.Errorf("key cannot contain spaces")
	}

	return e.setRaw(key, value)
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

const Seconds = 5

func (e *Engine) CompactFile() {
	for {
		time.Sleep(time.Duration(Seconds) * time.Second)
		fmt.Println("Compacting file...")
		e.mu.Lock()
		_, m := e.GetMapFromFile()
		e.mu.Unlock()

		e.fileMu.Lock()

		err := e.file.Truncate(0)
		if err != nil {
			fmt.Println(err)
			e.fileMu.Unlock()
			continue
		}

		for k, v := range m {
			e.setRaw(k, v)
		}

		e.file.Seek(0, 0)
		e.fileMu.Unlock()

	}
}



func (c *Engine) GetMapFromFile() ([]Item, map[string]string) {
	m := make(map[string]string)
	i := []Item{}

	_, err := c.file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return i, m
	}

	var totalBytesRead int64
	scanner := bufio.NewScanner(c.file)

	for scanner.Scan() {
		line := scanner.Text()
		offset := totalBytesRead
		parts := strings.Split(line, keyValueSeparator)
		if len(parts) >= 2 {
			m[parts[0]] = parts[1]
			i = append(i, Item{
				Key:    parts[0],
				Value:  parts[1],
				Offset: offset,
			})
		}
	}

	return i, m
}

func (e *Engine) setRaw(key string, value string) error {
	offset, err := e.saveToFile(key, value)
	if err != nil {
		return err
	}

	e.setKey(key, offset)
	return nil
}

func (e *Engine) setKey(key string, value int64) {
	e.data[key] = value
}

func (c *Engine) saveToFile(key string, value string) (int64, error) {
	offset, err := c.file.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return 0, err
	}

	_, err = c.file.WriteString(key + keyValueSeparator + value + "\n")
	if err != nil {
		fmt.Println("Error appending text:", err)
		return 0, err
	}

	return offset, nil
}

func (e *Engine) Restore() {
	e.mu.Lock()
	defer e.mu.Unlock()

	Items, _ := e.GetMapFromFile()

	for _, v := range Items {
		e.setKey(v.Key, v.Offset)
	}
}

func (c *Engine) Close() {
	c.file.Close()
}

func (e *Engine) Delete(key string) error {
	e.muDlete.Lock()
	defer e.muDlete.Unlock()

	_, err := e.fileDelete.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking delete file:", err)
		return err
	}

	_, err = e.fileDelete.WriteString(key + "\n")
	if err != nil {
		fmt.Println("Error appending delete text:", err)
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.data, key)

	return nil
}

const secondDelete = 5

func (e *Engine) DeleteFromFile() {
	for {
		time.Sleep(time.Duration(secondDelete) * time.Second)
		fmt.Println("Deleting from file...")
		e.muDlete.Lock()

		_, err := e.fileDelete.Seek(0, 0)
		if err != nil {
			fmt.Println("Error seeking delete file:", err)
			e.muDlete.Unlock()
			continue
		}

		scanner := bufio.NewScanner(e.fileDelete)

		content := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				content = append(content, line)
			}
		}

		err = e.deleteKeyFromFile(content)
		if err != nil {
			fmt.Println("Error deleting keys from file:", err)
			e.muDlete.Unlock()
			continue
		}

		err = e.fileDelete.Truncate(0)
		if err != nil {
			fmt.Println("Error truncating delete file:", err)
			e.muDlete.Unlock()
			continue
		}

		e.muDlete.Unlock()

	}
}

func (c *Engine) deleteKeyFromFile(keys []string) error {

	c.mu.Lock()
	defer c.mu.Unlock()
	c.fileMu.Lock() 
    defer c.fileMu.Unlock()
    

	_, err := c.file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return err
	}

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(c.file)
	for scanner.Scan() {
		l := scanner.Text()
		parts := strings.Split(l, keyValueSeparator)
		if len(parts) >= 2 {
			found := false
			for _, k := range keys {
				if parts[0] == k {
					found = true
					break
				}
			}

			if !found {
				buf.WriteString(l + "\n")
			}

		}

	}

	_, err = c.file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = c.file.Truncate(0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = buf.WriteTo(c.file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
