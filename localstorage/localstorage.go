package localstorage

import (
	"os"
	"sync"
)

var mtx sync.Mutex
var path string

func Init(programName string) {
	mtx.Lock()
	defer mtx.Unlock()

	if path != "" {
		return
	}

	homeDir := homeDirectory()
	path = homeDir + "/." + programName

	os.MkdirAll(path, 0600)
}

func Path() string {
	mtx.Lock()
	defer mtx.Unlock()
	return path
}

func Write(fileName string, data []byte) error {
	mtx.Lock()
	defer mtx.Unlock()
	filePath := path + "/" + fileName
	err := os.WriteFile(filePath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func Read(fileName string) ([]byte, error) {
	mtx.Lock()
	defer mtx.Unlock()
	filePath := path + "/" + fileName
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Exists(fileName string) bool {
	mtx.Lock()
	defer mtx.Unlock()
	filePath := path + "/" + fileName
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
