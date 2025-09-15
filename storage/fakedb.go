// Package storage serves a fake DB right now
package storage

import (
	"fmt"
	"os"
	"sync"
)

var (
	Users          sync.Map
	requestCounter = make(map[string]int)
)

func UpdateRequestCounter(key string) {
	if val, ok := requestCounter[key]; ok {
		val++
		requestCounter[key] = val
	} else {
		requestCounter[key] = 1
	}

	fileContent := ""
	for key, val := range requestCounter {
		fileContent += fmt.Sprintf("%s: %d\n", key, val)
	}

	file, err := os.Create("storage/visitors.txt")
	if err != nil {
		return
	}

	file.Write([]byte(fileContent))
}
