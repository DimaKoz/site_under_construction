package app

import (
	"io/ioutil"
	"sync"
)

var cacheByte = make(map[string][]byte)
var cacheByteMutex sync.Mutex

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

func RemoveKey(key string) {
	cacheByteMutex.Lock()
	defer cacheByteMutex.Unlock()

	delete(cacheByte, key)
}

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
