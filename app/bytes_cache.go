package app

import (
	"io/ioutil"
	"sync"
)

var cacheByteMutex sync.Mutex
var cacheByte = make(map[string][]byte)

// AddKeyAndPath appends to the internal cache with 'key' a file with 'path'
// It needs for some tests
func AddKeyAndPath(key string, path string) error {
	cacheByteMutex.Lock()
	defer cacheByteMutex.Unlock()

	bytesFromCache := cacheByte[key]
	if bytesFromCache == nil {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		cacheByte[key] = data
	}
	return nil
}

// RemoveKey removes a value from a cache by 'key'
func RemoveKey(key string) {
	cacheByteMutex.Lock()
	defer cacheByteMutex.Unlock()

	delete(cacheByte, key)
}

// GetBytes reads data from file puts it to a cache by 'path' key
func GetBytes(fileName string) (*[]byte, error) {
	cacheByteMutex.Lock()
	defer cacheByteMutex.Unlock()

	bytesFromCache := cacheByte[fileName]
	if bytesFromCache == nil {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		cacheByte[fileName] = data

	}
	bytesFromCache = cacheByte[fileName]
	result := make([]byte, len(bytesFromCache))
	copy(result, bytesFromCache)
	return &result, nil
}
