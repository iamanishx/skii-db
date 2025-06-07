package engine

import (
	"fmt"
	"os"
	"testing"
)

// func Test_SetGetKeyValue(t *testing.T) {
// 	e, _ := NewEngine()
// 	e.Set("foo", "bar")
// 	value, err := e.Get("foo")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if value != "bar" {
// 		t.Error("value should be bar")
// 	}

// 	_, err = e.Get("notfound")
// 	if err == nil {
// 		t.Error("should return error")
// 	}

// }

// func (c *Engine) GetFileContent(f *os.File) []string {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	_, err := f.Seek(0, 0)
// 	if err != nil {
// 		fmt.Println(err)
// 		return []string{}
// 	}

// 	scanner := bufio.NewScanner(f)

// 	var content []string
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		content = append(content, line)
// 	}
// 	return content
// }

// func TestEngine_Compact(t *testing.T) {
//     os.Remove("data.txt")
// 	v1 := "latestvalue1"
// 	v2 := "latestvalue2"
// 	e, _ := NewEngine()
// 	e.Set("key1", "value1")
// 	e.Set("key2", "value2")
// 	e.Set("key1", v1)
// 	e.Set("key2", v2)
// 	e.Set("key3", "value3")

// 	go e.CompactFile()

// 	time.Sleep((Seconds + 3) * time.Second)
// 	if len(e.GetFileContent(e.file)) != 3 {
// 		t.Errorf("Expected %d, but got %d", 3, len(e.GetFileContent(e.file)))
// 	}

// }

func TestEngine_Restore(t *testing.T) {
	os.Remove("data.txt")
	e, _ := NewEngine()

	e.Set("key1_restore", "value1")
	e.Set("key2_restore", "value2")

	e.Close()
	e, _ = NewEngine()
	e.Restore()

	k, _ := e.Get("key1_restore")

	fmt.Println("Restored key1:", k)
	if k != "value1" {
		t.Errorf("Expected value1, but got %s", k)
	}

}
