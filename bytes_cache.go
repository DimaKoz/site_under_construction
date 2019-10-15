package main

import (
	"io/ioutil"
	"sync"
)

var cacheByte = make(map[string][]byte)
var cacheByteMutex sync.Mutex


func getBytes(fileName string) (*[]byte, error) {
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

